variable "region" {
  description = "Region for AWS resources"
  default     = "eu-west-2"
}

variable "bucket" {
  description = "Name of the S3 bucket"
}

variable "repository" {
  description = "Repository name. This will be used as a folder in s3 to store policies and config files"
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

variable lambda_zip_file {
  description = "Zip file containing the lambda function"
}

variable "lambda_function_name" {
  description = "Lambda function name"
}

variable "github_token" {
  description = "Token to call GitHub REST APIs"
}

variable "gitlab_token" {
  description = "Token to call GitLab REST APIs"
}

variable "slack_token" {
  description = "Token to call Slack APIs for notifications"
}

variable "event_schedule_rate" {
  description = "Rate for the event, in the form of 'rate(value unit)'. value: a positive number, unit: minute | minutes | hour | hours | day | days"
  default = "rate(12 hours)"
}