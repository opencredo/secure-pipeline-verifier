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

# Provision config folders for the repositories in the S3 bucket
module "repositories" {
  source = "./modules/repository"
  for_each = { for repo in var.repo_list : repo.path => repo }
  source_dir  = each.key
  bucket      = aws_s3_bucket.secure_pipeline.bucket
  lambda_arn  = aws_lambda_function.check_policies.arn
  lambda_name = aws_lambda_function.check_policies.function_name
  last_run    = timestamp()
  parameter_prefix = var.parameter_prefix
  repo_token  = each.value.repo_token
  region = var.region
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

resource "aws_lambda_function" "check_policies" {
  filename         = var.lambda_zip_file
  function_name    = var.lambda_function_name
  source_code_hash = filebase64sha256(var.lambda_zip_file)
  timeout          = var.lambda_timeout
  role             = aws_iam_role.lambda.arn
  handler          = "main"
  runtime          = "go1.x"
}
