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
  bucket = aws_s3_bucket.secure-pipeline.bucket
  key    = "${var.repository}/config.yaml"
  source = var.config_file
}

resource "aws_s3_bucket_object" "trusted_data_file" {
  bucket = aws_s3_bucket.secure-pipeline.bucket
  key    = "${var.repository}/trusted_data.json"
  source = var.trusted_data_file
}

resource "aws_s3_bucket_object" "policies" {
  bucket   = aws_s3_bucket.secure-pipeline.bucket
  for_each = fileset(var.policies_dir, "*.rego")
  key      = "${var.repository}/policies/${each.value}"
  source   = "${var.policies_dir}/${each.value}"
}

resource "aws_ssm_parameter" "last_run" {
  description = "Last run of Secure Pipeline. Format: 'YYYY-MM-ddTHH:mm:ss.SSSZ'."
  name  = "lambda/secure_pipeline/last_run"
  type  = "String"
  value = "" // a single whitespace
  overwrite = false
}
