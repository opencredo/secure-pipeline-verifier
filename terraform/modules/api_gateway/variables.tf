variable "invoke_arn" {
  type        = string
  description = "Lambda arn for AGW integration"
}

variable "path_part" {
  type        = string
  description = "The last path segment of this API resource."
}
