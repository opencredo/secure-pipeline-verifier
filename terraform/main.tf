terraform {
  experiments = [module_variable_optional_attrs]
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
  }
}

provider "aws" {
  region = var.region
}

resource "aws_s3_bucket" "secure_pipeline" {
  bucket = var.bucket
  acl    = "private"
}

# Provision resources for the repositories to be audited by the application
module "repositories" {
  source              = "./modules/repository"
  for_each            = { for repo in var.repo_list : repo.path => repo }
  source_dir          = each.key
  event_schedule_rate = coalesce(each.value.event_schedule_rate, var.event_schedule_rate)
  repo_token          = each.value.repo_token
  bucket              = aws_s3_bucket.secure_pipeline.bucket
  lambda_arn          = aws_lambda_function.check_policies.arn
  lambda_name         = aws_lambda_function.check_policies.function_name
  last_run            = coalesce(var.last_run, timestamp())
  parameter_prefix    = var.parameter_prefix
  region              = var.region
}

resource "aws_ssm_parameter" "slack_token" {
  description = "A token to authenticate with a repository."
  name        = "${var.parameter_prefix}/SLACK_TOKEN"
  type        = "SecureString"
  value       = var.slack_token
}

resource "aws_cloudwatch_log_group" "lambda" {
  name = "/aws/lambda/${var.lambda_function_name}"
}

resource "aws_cloudwatch_log_group" "cw_chatops" {
  count = var.lambda_chatops_zip_file != null ? 1 : 0
  name  = "/aws/lambda/${var.lambda_chatops_name}"
}

resource "aws_cloudwatch_log_stream" "lambda_chatops" {
  count          = var.lambda_chatops_zip_file != null ? 1 : 0
  log_group_name = aws_cloudwatch_log_group.cw_chatops[0].name
  name           = "lambda-stream-chatops"
}

resource "aws_cloudwatch_log_stream" "lambda" {
  log_group_name = aws_cloudwatch_log_group.lambda.name
  name           = "lambda-stream"
}

data "aws_caller_identity" "current" {}

resource "aws_iam_role" "lambda" {
  name = "SecurePipelineLambdaAccess"
  assume_role_policy = jsonencode({
    Version = "2008-10-17",
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          "Service" : "lambda.amazonaws.com"
        }
      },
    ]
  })
  inline_policy {
    name = "LambdaAccessToServices"
    policy = jsonencode({
      "Statement" : [
        {
          "Effect" : "Allow",
          "Action" : [
            "s3:GetObject",
            "s3:ListBucket",
            "logs:CreateLogStream",
            "logs:CreateLogGroup",
            "logs:PutLogEvents",
          ],
          "Resource" : [
            "arn:aws:logs:${var.region}:${data.aws_caller_identity.current.account_id}:log-group:${aws_cloudwatch_log_group.lambda.name}:*",
            "arn:aws:s3:::${aws_s3_bucket.secure_pipeline.bucket}",
            "arn:aws:s3:::${aws_s3_bucket.secure_pipeline.bucket}/*",
          ]
        },
        {
          "Effect" : "Allow",
          "Action" : [
            "ssm:PutParameter",
            "ssm:GetParameter",
          ],
          "Resource" : [
            "arn:aws:ssm:${var.region}:${data.aws_caller_identity.current.account_id}:parameter${var.parameter_prefix}/*",
          ]
        }
      ]
    })
  }
  depends_on = [
    aws_s3_bucket.secure_pipeline,
    aws_cloudwatch_log_group.lambda,
  ]
}

resource "aws_iam_role" "call_lambda" {
  count          = var.lambda_chatops_zip_file != null ? 1 : 0
  name  = "ChatOpsCallLambda"
  assume_role_policy = jsonencode({
    Version = "2008-10-17",
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          "Service" : "lambda.amazonaws.com"
        }
      },
    ]
  })
  inline_policy {
    name = "LambdaCallLambda"
    policy = jsonencode({
      "Statement" : [
        {
          "Effect" : "Allow",
          "Action" : [
            "logs:CreateLogStream",
            "logs:CreateLogGroup",
            "logs:PutLogEvents",
          ],
          "Resource" : [
            "arn:aws:logs:${var.region}:${data.aws_caller_identity.current.account_id}:log-group:${aws_cloudwatch_log_group.cw_chatops[0].name}:*",
          ]
        },
        {
          "Effect" : "Allow",
          "Action" : [
            "lambda:InvokeFunction",
          ],
          "Resource" : [
            "arn:aws:lambda:${var.region}:${data.aws_caller_identity.current.account_id}:function:${var.lambda_function_name}",
          ]
        }
      ]
    })
  }

}

resource "aws_lambda_function" "check_policies" {
  filename         = var.lambda_zip_file
  function_name    = var.lambda_function_name
  source_code_hash = filebase64sha256(var.lambda_zip_file)
  timeout          = var.lambda_timeout
  role             = aws_iam_role.lambda.arn
  handler          = "main"
  runtime          = "go1.x"
}

resource "aws_lambda_function" "chatops" {
  count            = var.lambda_chatops_zip_file != null ? 1 : 0
  filename         = var.lambda_chatops_zip_file
  function_name    = var.lambda_chatops_name
  source_code_hash = filebase64sha256(var.lambda_chatops_zip_file)
  timeout          = var.lambda_timeout
  role             = aws_iam_role.call_lambda[0].arn
  handler          = "main"
  runtime          = "go1.x"

  environment {
    variables = {
      TARGET_LAMBDA = aws_lambda_function.check_policies.function_name
    }
  }

}

# API Gateway
resource "aws_api_gateway_rest_api" "api" {
  name = "secure-pipeline-api"
  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

module "api_gateway_lambda" {
  source           = "./modules/api_gateway"
  path_part        = "audit"
  api_id           = aws_api_gateway_rest_api.api.id
  root_resource_id = aws_api_gateway_rest_api.api.root_resource_id
  account_id       = data.aws_caller_identity.current.account_id
  function_name    = var.lambda_function_name
  invoke_arn       = aws_lambda_function.check_policies.invoke_arn
  depends_on = [
    aws_lambda_function.check_policies
  ]
}

module "api_gateway_lambda_chatops" {
  count                = var.lambda_chatops_zip_file != null ? 1 : 0
  source               = "./modules/api_gateway"
  path_part            = "chatops"
  api_id               = aws_api_gateway_rest_api.api.id
  root_resource_id     = aws_api_gateway_rest_api.api.root_resource_id
  account_id           = data.aws_caller_identity.current.account_id
  function_name        = var.lambda_chatops_name
  invoke_arn           = aws_lambda_function.chatops[0].invoke_arn
  passthrough_behavior = "WHEN_NO_TEMPLATES"
  urlencoded_tmpl      = <<-EOT
                {
                  "body" : $input.json('$')
                }
    EOT
  depends_on = [
    aws_lambda_function.chatops
  ]

}

resource "aws_api_gateway_deployment" "api_deploy" {
  rest_api_id = aws_api_gateway_rest_api.api.id

  triggers = {
    redeployment = sha1(
      join(",",
        tolist(
          flatten([
            jsonencode(module.api_gateway_lambda.lambda_integration),
            var.lambda_chatops_zip_file != null ? [jsonencode(module.api_gateway_lambda_chatops[0].lambda_integration)] : []
    ]))))
  }

  lifecycle {
    create_before_destroy = true
  }

}

resource "aws_api_gateway_stage" "v1" {
  deployment_id = aws_api_gateway_deployment.api_deploy.id
  rest_api_id   = aws_api_gateway_rest_api.api.id
  stage_name    = "v1"
}

output "api_url" {
  value = aws_api_gateway_stage.v1.invoke_url
}
