locals {
  # Get information about the repository
  config_file       = "${var.source_dir}/config.yaml"
  trusted_data_file = "${var.source_dir}/trusted-data.yaml"
  policies_dir      = "${var.source_dir}/policies"
  # Open the config file to access the 'project' block.
  project        = yamldecode(file(local.config_file))["project"]
  platform       = local.project["platform"]
  repo_name      = local.project["repo"]
  parameter_path = "${var.parameter_prefix}/${local.project["platform"]}/${local.project["owner"]}/${local.repo_name}"
}

resource "aws_s3_bucket_object" "config_file" {
  for_each    = fileset(var.source_dir, "*.yaml")
  key         = "${local.repo_name}/${basename(local.config_file)}"
  bucket      = var.bucket
  source      = local.config_file
  source_hash = filemd5(local.config_file)
}

resource "aws_s3_bucket_object" "trusted_data_file" {
  bucket      = var.bucket
  key         = "${local.repo_name}/${basename(local.trusted_data_file)}"
  source      = local.trusted_data_file
  source_hash = filemd5(local.trusted_data_file)
}

resource "aws_s3_bucket_object" "policies" {
  bucket      = var.bucket
  for_each    = fileset(local.policies_dir, "*.rego")
  key         = "${local.repo_name}/${basename(local.policies_dir)}/${each.value}"
  source      = "${local.policies_dir}/${each.value}"
  source_hash = filemd5("${local.policies_dir}/${each.value}")
}

resource "aws_ssm_parameter" "last_run" {
  description = "Last run of Secure Pipeline. Format: 'YYYY-MM-DD'T'hh:mm:ssZ'."
  name        = "${local.parameter_path}/last_run"
  type        = "String"
  value       = var.last_run
  lifecycle {
    # Fill the value when the resource is created for the first time. Later it might be changed outside of Terraform.
    ignore_changes = [
      value,
    ]
  }
}

resource "aws_ssm_parameter" "repo_token" {
  description = "A token to authenticate with a repository."
  name        = "${local.parameter_path}/REPO_TOKEN"
  type        = "SecureString"
  value       = var.repo_token
}

resource "aws_cloudwatch_event_rule" "trigger_lambda_event_rule" {
  name                = "trigger_lambda_event_rule_${local.platform}_${local.repo_name}"
  description         = "Fires Lambda execution"
  schedule_expression = var.event_schedule_rate
}

resource "aws_cloudwatch_event_target" "check_policies_event_target" {
  rule      = aws_cloudwatch_event_rule.trigger_lambda_event_rule.name
  target_id = "check_policies"
  arn       = var.lambda_arn
  input = jsonencode({
    "region" : var.region,
    "bucket" : var.bucket,
    "configPath" : local.project["repo"]
  })
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_check_policies" {
  statement_id  = "AllowExecutionFromCloudWatchFor_${local.platform}${local.repo_name}"
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.trigger_lambda_event_rule.arn
}
