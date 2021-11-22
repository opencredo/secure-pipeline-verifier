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

locals {
  project      = yamldecode(file(var.config_file))["project"]
  repo_name    =  local.project["repo"]
  parameterPrefix = "/Lambda/SecurePipelines"
  parameterPath = "${local.parameterPrefix}/${local.project["platform"]}/${local.project["owner"]}/${local.repo_name}"
}


resource "aws_s3_bucket" "secure_pipeline" {
  bucket = var.bucket
  acl    = "private"
}

resource "aws_s3_bucket_object" "config_file" {
  bucket      = aws_s3_bucket.secure_pipeline.bucket
  key         = "${local.repo_name}/config.yaml"
  source      = var.config_file
  source_hash = filemd5(var.config_file)
  depends_on  = [aws_s3_bucket.secure_pipeline]
}

resource "aws_s3_bucket_object" "trusted_data_file" {
  bucket      = aws_s3_bucket.secure_pipeline.bucket
  key         = "${local.repo_name}/trusted-data.yaml"
  source      = var.trusted_data_file
  source_hash = filemd5(var.trusted_data_file)
  depends_on  = [aws_s3_bucket.secure_pipeline]
}

resource "aws_s3_bucket_object" "policies" {
  bucket      = aws_s3_bucket.secure_pipeline.bucket
  for_each    = fileset(var.policies_dir, "*.rego")
  key         = "${local.repo_name}/policies/${each.value}"
  source      = "${var.policies_dir}/${each.value}"
  source_hash = filemd5("${var.policies_dir}/${each.value}")
  depends_on  = [aws_s3_bucket.secure_pipeline]
}

resource "aws_ssm_parameter" "last_run" {
  description = "Last run of Secure Pipeline. Format: 'YYYY-MM-DD'T'hh:mm:ssZ'."
  name        = "${local.parameterPath}/last_run"
  type        = "String"
  value = var.last_run
  lifecycle {
    # Fill the value when the resource is created for the first time. Later it might be changed outside of Terraform.
    ignore_changes = [
      value,
    ]
  }
}

resource "aws_ssm_parameter" "repo_token" {
  description = "A token to authenticate with a repository."
  name        = "${local.parameterPath}/repo_token"
  type        = "SecureString"
  value       = var.repo_token
}


resource "aws_ssm_parameter" "slack_token" {
  description = "A token to authenticate with a repository."
  name        = "${local.parameterPrefix}/slack_token"
  type        = "SecureString"
  value       = var.slack_token
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
            "arn:aws:ssm:${var.region}:${data.aws_caller_identity.current.account_id}:parameter${local.parameterPrefix}/*",
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
    "configPath" : local.project["repo"]
  })
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_check_policies" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.check_policies.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.trigger_lambda_event_rule.arn
}
