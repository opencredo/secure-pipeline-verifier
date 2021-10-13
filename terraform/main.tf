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

resource "aws_lambda_function" "check_policies" {
  filename = var.lambda_zip_file
  function_name = var.lambda_function_name
  role = aws_iam_role.lambda.arn
  handler = "main"
}

resource "aws_cloudwatch_event_rule" "trigger_lambda_event_rule" {
  name = "trigger_lambda_event_rule"
  description = "Fires Lambda execution"
  schedule_expression = var.event_schedule_rate
}

resource "aws_cloudwatch_event_target" "check_policies_event_target" {
  rule = aws_cloudwatch_event_rule.trigger_lambda_event_rule.name
  target_id = "check_policies"
  arn = aws_lambda_function.check_policies.arn
  input = <<JSON
  "{
    "region": "${var.region}",
    "bucket": "${aws_s3_bucket.secure_pipeline.bucket}",
    "configPath": ""
  }"
  JSON
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_check_policies" {
  statement_id = "AllowExecutionFromCloudWatch"
  action = "lambda:InvokeFunction"
  function_name = aws_lambda_function.check_policies.function_name
  principal = "events.amazonaws.com"
  source_arn = aws_cloudwatch_event_rule.trigger_lambda_event_rule.arn
}