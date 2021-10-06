variable "region" {
  description = "Region for AWS resources"
  default = "eu-west-2"
}

variable "bucket" {
  description = "Name of the S3 bucket"
}

variable "platform" {
  description = "Choose which repository to audit. For example: github"
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
