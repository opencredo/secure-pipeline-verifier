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
  key    = "${var.repository}/config.yaml"
  source = var.config_file
}

resource "aws_s3_bucket_object" "trusted_data_file" {
  bucket = aws_s3_bucket.secure_pipeline.bucket
  key    = "${var.repository}/trusted_data.json"
  source = var.trusted_data_file
}

resource "aws_s3_bucket_object" "policies" {
  bucket   = aws_s3_bucket.secure_pipeline.bucket
  for_each = fileset(var.policies_dir, "*.rego")
  key      = "${var.repository}/policies/${each.value}"
  source   = "${var.policies_dir}/${each.value}"
}

resource "aws_ssm_parameter" "last_run" {
  description = "Last run of Secure Pipeline. Format: 'YYYY-MM-DD'T'hh:mm:ssZ'."
  name  = "/Lambda/SecurePipelines/last_run"
  type  = "String"
  # If the value doesn't exist then the last run will be the deployment time of this resource.
  value = timestamp()
  lifecycle {
    # Fill the value when the resource is created for the first time. Later it might be changed outside of Terraform.
    ignore_changes = [
      value,
  ]
  }
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
    "region": "${aws_s3_bucket.secure_pipeline.region}",
    "bucket": "${aws_s3_bucket.secure_pipeline.bucket}",
    "configPath": "${var.repository}"
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
