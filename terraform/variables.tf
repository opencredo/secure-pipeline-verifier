variable "region" {
  description = "Region for AWS resources"
  default = "eu-west-2"
}

variable lambda_zip_file {
  description = "Zip file containing the lambda function"
}

variable "lambda_function_name" {
  description = "Lambda function name"
}

variable "event_schedule_rate" {
  description = "Rate for the event, in the form of 'rate(value unit)'. value: a positive number, unit: minute | minutes | hour | hours | day | days"
  default = "rate(12 hours)"
}