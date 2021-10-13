terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
  }
}

provider "aws" {
  region = var.region
}

resource "aws_s3_bucket" "secure_pipeline" {
  bucket = var.bucket
  acl    = "private"
}

resource "aws_s3_bucket_object" "config_file" {
  bucket = aws_s3_bucket.secure_pipeline.bucket
  key    = "${var.repository}/config.yaml"
  source = var.config_file
}

resource "aws_s3_bucket_object" "trusted_data_file" {
  bucket = aws_s3_bucket.secure_pipeline.bucket
  key    = "${var.repository}/trusted_data.json"
  source = var.trusted_data_file
}

resource "aws_s3_bucket_object" "policies" {
  bucket   = aws_s3_bucket.secure_pipeline.bucket
  for_each = fileset(var.policies_dir, "*.rego")
  key      = "${var.repository}/policies/${each.value}"
  source   = "${var.policies_dir}/${each.value}"
}

resource "aws_ssm_parameter" "last_run" {
  description = "Last run of Secure Pipeline. Format: 'YYYY-MM-DD'T'hh:mm:ssZ'."
  name  = "/Lambda/SecurePipelines/last_run"
  type  = "String"
  # If the value doesn't exist then the last run will be the deployment time of this resource.
  value = timestamp()
  lifecycle {
    # Fill the value when the resource is created for the first time. Later it might be changed outside of Terraform.
    ignore_changes = [
      value,
  ]
  }
}
