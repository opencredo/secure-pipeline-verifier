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
  bucket = aws_s3_bucket.secure_pipeline.bucket
  key    = "${var.platform}/config.yaml"
  source = var.config_file
  depends_on = [aws_s3_bucket.secure_pipeline]
}

resource "aws_s3_bucket_object" "trusted_data_file" {
  bucket = aws_s3_bucket.secure_pipeline.bucket
  key    = "${var.platform}/trusted_data.json"
  source = var.trusted_data_file
  depends_on = [aws_s3_bucket.secure_pipeline]
  lifecycle {
    prevent_destroy = true
  }
}

resource "aws_s3_bucket_object" "policies" {
  bucket   = aws_s3_bucket.secure_pipeline.bucket
  key      = "${var.platform}/policies/${each.value}"
  source   = "${var.policies_dir}/${each.value}"
  for_each = fileset(var.policies_dir, "*.rego")
  depends_on = [aws_s3_bucket.secure_pipeline]
}

resource "aws_cloudwatch_log_group" "lambda" {
  name = "lambda-logs"
}

resource "aws_cloudwatch_log_stream" "lambda" {
  log_group_name = aws_cloudwatch_log_group.lambda.name
  name           = "lambda-stream"
  depends_on = [
    aws_cloudwatch_log_group.lambda
  ]
}

data "aws_caller_identity" "current" {}


resource "aws_ssm_parameter" "last_run" {
  description = "Last run of Secure Pipeline. Format: 'YYYY-MM-ddTHH:mm:ss.SSSZ'."
  name  = "lambda/secure_pipeline/last_run"
  type  = "String"
  value = " " // a single whitespace
  overwrite = false
}


resource "aws_iam_role" "lambda" {
  name = "SecurePipelineLambdaAccess"
  assume_role_policy = jsonencode({
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          "Service": "lambda.amazonaws.com"
        }
      },
    ]
  })
  inline_policy {
    name = "LambdaAccessToServices"
    policy = jsonencode({
      "Statement": [
        {
          "Effect": "Allow",
          "Action": [
            "s3:GetObject",
            "logs:CreateLogStream",
            "logs:CreateLogGroup",
            "logs:PutLogEvents",
          ],
          "Resource": [
            "arn:aws:logs:${var.region}:${data.aws_caller_identity.current.account_id}:log-group:${aws_cloudwatch_log_group.lambda.name}",
            "arn:aws:logs:${var.region}:${data.aws_caller_identity.current.account_id}:log-group:${aws_cloudwatch_log_group.lambda.name}:log-stream:${aws_cloudwatch_log_stream.lambda.name}",
            "arn:aws:s3:::${aws_s3_bucket.secure_pipeline.bucket}/*",
          ]
        },
        {
          "Effect": "Allow",
          "Action": [
            "ssm:PutParameter",
            "ssm:GetParameter",
          ],
          "Resource": [
            "arn:aws:ssm:${var.region}:${data.aws_caller_identity.current.account_id}:parameter/${aws_ssm_parameter.last_run.name}",
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

