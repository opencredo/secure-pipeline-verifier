variable "region" {
  description = "Region for AWS resources"
  default     = "eu-west-2"
}

variable "source_dir" {
  type = string
  description = "A directory containing config (.yaml) and policy (.rego) files for the repository"
}

variable "bucket" {
  description = "Name of the S3 bucket"
}

variable "repo_token" {
  description = "Token to call a Version Control REST APIs"
  sensitive   = true
}

variable "lambda_arn" {
  description = "Lambda resource name"
}

variable "lambda_name" {
  description = "Name of the lambda function"
}

variable "parameter_prefix" {
  description = "A path in the parameter store to save the configs for this repository"
}

variable "event_schedule_rate" {
  description = "Rate for the event, in the form of 'rate(value unit)'. value: a positive number, unit: minute | minutes | hour | hours | day | days"
  default     = "rate(12 hours)"
}

variable "last_run" {
  description = "Last run of Secure Pipeline service. If first run, set its value to a date in the past where you want to start verifying policies. Format: 'YYYY-MM-DD'T'hh:mm:ssZ'. "
}