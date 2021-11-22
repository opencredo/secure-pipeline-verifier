variable "region" {
  description = "Region for AWS resources"
  default     = "eu-west-2"
}

variable "bucket" {
  description = "Name of the S3 bucket"
}

variable "policies_dir" {
  description = "Path to a directory with policy files (*.rego)"
}

variable "config_file" {
  description = "Config file for Secure Pipeline"
}

variable "trusted_data_file" {
  description = "JSON file for policies"
}

variable "lambda_zip_file" {
  description = "Zip file containing the lambda function"
}

variable "lambda_function_name" {
  description = "Lambda function name"
  default     = "secure_pipeline"
}

variable "lambda_timeout" {
  description = "Amount of time your Lambda Function has to run in seconds."
  default     = 10
}

variable "repo_token" {
  description = "Token to call a Version Control REST APIs"
}

variable "slack_token" {
  description = "Token to call Slack APIs for notifications"
}

variable "event_schedule_rate" {
  description = "Rate for the event, in the form of 'rate(value unit)'. value: a positive number, unit: minute | minutes | hour | hours | day | days"
  default     = "rate(12 hours)"
}

variable "last_run" {
  description = "Last run of Secure Pipeline service. If first run, set its value to a date in the past where you want to start verifying policies. Format: 'YYYY-MM-DD'T'hh:mm:ssZ'. "
}