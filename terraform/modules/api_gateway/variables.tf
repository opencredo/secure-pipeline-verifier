variable "function_name" {
  type        = string
  description = "Name of the lambda function"
}

variable "account_id" {
  type        = string
  description = "Account ID that is authorised to use Terraform"
}

variable "invoke_arn" {
  type        = string
  description = "Lambda arn for AGW integration"
}

variable "path_part" {
  type        = string
  description = "The last path segment of this API resource."
}

variable "region" {
  type        = string
  description = "Region for AWS resources"
  default     = "eu-west-2"
}
