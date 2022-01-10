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

variable "api_id" {
  type        = string
  description = "ID of the API gateway."
}

variable "root_resource_id" {
  type        = string
  description = "ID of the root resource."
}

variable "urlencoded_tmpl" {
  type        = string
  description = "A Template for application/x-www-form-urlencoded"
  default     = ""
}

variable "passthrough_behavior" {
  type        = string
  description = "A passthrough behavior for the lambda integration. Default: WHEN_NO_MATCH."
  default     = "WHEN_NO_MATCH"
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
