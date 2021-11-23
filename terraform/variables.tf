variable "region" {
  description = "Region for AWS resources"
  default     = "eu-west-2"
}

variable "repo_list" {
  description = <<EOF
  A list of maps.
     - path: path to a directory containing config files for a specific repository,
     - repo_token: Token to call a Version Control REST APIs
  EOF
  type = list(object({
    path       = string
    repo_token = string
  }))
}

variable "bucket" {
  description = "Name of the S3 bucket"
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

variable "parameter_prefix" {
  description = "A path in the parameter store to save the configs for this repository"
  default     = "/Lambda/SecurePipelines"
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