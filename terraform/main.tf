terraform {
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

resource "aws_s3_bucket_object" "config_file" {
  bucket      = aws_s3_bucket.secure_pipeline.bucket
  key         = "${var.repository}/config.yaml"
  source      = var.config_file
  source_hash = filemd5(var.config_file)
  depends_on  = [aws_s3_bucket.secure_pipeline]
}

resource "aws_s3_bucket_object" "trusted_data_file" {
  bucket      = aws_s3_bucket.secure_pipeline.bucket
  key         = "${var.repository}/trusted-data.yaml"
  source      = var.trusted_data_file
  source_hash = filemd5(var.trusted_data_file)
  depends_on  = [aws_s3_bucket.secure_pipeline]
}

resource "aws_s3_bucket_object" "policies" {
  bucket      = aws_s3_bucket.secure_pipeline.bucket
  for_each    = fileset(var.policies_dir, "*.rego")
  key         = "${var.repository}/policies/${each.value}"
  source      = "${var.policies_dir}/${each.value}"
  source_hash = filemd5("${var.policies_dir}/${each.value}")
  depends_on  = [aws_s3_bucket.secure_pipeline]
}

resource "aws_ssm_parameter" "last_run" {
  description = "Last run of Secure Pipeline. Format: 'YYYY-MM-DD'T'hh:mm:ssZ'."
  name        = "/Lambda/SecurePipelines/last_run"
  type        = "String"
  value = var.last_run
  lifecycle {
    # Fill the value when the resource is created for the first time. Later it might be changed outside of Terraform.
    ignore_changes = [
      value,
    ]
  }
}

resource "aws_cloudwatch_log_group" "lambda" {
  name = "/aws/lambda/${var.lambda_function_name}"
}

resource "aws_cloudwatch_log_stream" "lambda" {
  log_group_name = aws_cloudwatch_log_group.lambda.name
  name           = "lambda-stream"
}

data "aws_caller_identity" "current" {}

resource "aws_iam_role" "lambda" {
  name = "SecurePipelineLambdaAccess"
  assume_role_policy = jsonencode({
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
            "arn:aws:ssm:${var.region}:${data.aws_caller_identity.current.account_id}:parameter${aws_ssm_parameter.last_run.name}",
          ]
        }
      ]
    })
  }
  depends_on = [
    aws_s3_bucket.secure_pipeline,
    aws_cloudwatch_log_group.lambda,
    aws_ssm_parameter.last_run
  ]
}

resource "aws_lambda_function" "check_policies" {
  filename         = var.lambda_zip_file
  function_name    = var.lambda_function_name
  source_code_hash = filebase64sha256(var.lambda_zip_file)
  timeout          = var.lambda_timeout
  role             = aws_iam_role.lambda.arn
  handler          = "main"
  runtime          = "go1.x"

  environment {
    variables = {
      GITHUB_TOKEN = sensitive(var.github_token)
      GITLAB_TOKEN = sensitive(var.gitlab_token)
      SLACK_TOKEN  = sensitive(var.slack_token)
    }
  }
}

resource "aws_cloudwatch_event_rule" "trigger_lambda_event_rule" {
  name                = "trigger_lambda_event_rule"
  description         = "Fires Lambda execution"
  schedule_expression = var.event_schedule_rate
}

resource "aws_cloudwatch_event_target" "check_policies_event_target" {
  rule      = aws_cloudwatch_event_rule.trigger_lambda_event_rule.name
  target_id = "check_policies"
  arn       = aws_lambda_function.check_policies.arn
  input = jsonencode({
    "region" : aws_s3_bucket.secure_pipeline.region,
    "bucket" : aws_s3_bucket.secure_pipeline.bucket,
    "configPath" : var.repository
  })
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_check_policies" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.check_policies.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.trigger_lambda_event_rule.arn
}

resource "aws_lambda_permission" "allow_api_gw_to_call_check_policies" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.check_policies.function_name
  principal     = "apigateway.amazonaws.com"
  # More: http://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-control-access-using-iam-policies-to-invoke-api.html
  source_arn = "arn:aws:execute-api:${var.region}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.api.id}/*/${aws_api_gateway_method.method.http_method}${aws_api_gateway_resource.resource.path}"
}

# API Gateway
resource "aws_api_gateway_rest_api" "api" {
  name = "myapi"
  endpoint_configuration {
    types = ["REGIONAL"]
  }
}


resource "aws_api_gateway_resource" "resource" {
  path_part   = "audit"
  parent_id   = aws_api_gateway_rest_api.api.root_resource_id
  rest_api_id = aws_api_gateway_rest_api.api.id
}

resource "aws_api_gateway_method" "method" {
  rest_api_id   = aws_api_gateway_rest_api.api.id
  resource_id   = aws_api_gateway_resource.resource.id
  http_method   = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_method_response" "response_200" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.resource.id
  http_method = aws_api_gateway_method.method.http_method
  status_code = "200"
}

resource "aws_api_gateway_integration_response" "integration_response" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.resource.id
  http_method = aws_api_gateway_method.method.http_method
  status_code = aws_api_gateway_method_response.response_200.status_code
}

resource "aws_api_gateway_integration" "integration" {
  rest_api_id             = aws_api_gateway_rest_api.api.id
  resource_id             = aws_api_gateway_resource.resource.id
  http_method             = aws_api_gateway_method.method.http_method
  integration_http_method = aws_api_gateway_method.method.http_method
  type                    = "AWS"
  uri                     = aws_lambda_function.check_policies.invoke_arn
}

resource "aws_api_gateway_deployment" "api_deploy" {
  rest_api_id = aws_api_gateway_rest_api.api.id

  triggers = {
    redeployment = sha1(jsonencode(aws_api_gateway_rest_api.api.body))
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_stage" "staging" {
  deployment_id = aws_api_gateway_deployment.api_deploy.id
  rest_api_id   = aws_api_gateway_rest_api.api.id
  stage_name    = "staging"
}
