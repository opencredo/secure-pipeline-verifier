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

resource "aws_s3_bucket" "secure-pipeline" {
  bucket = var.bucket
  acl    = "private"
}

resource "aws_s3_bucket_object" "config_file" {
  bucket = aws_s3_bucket.secure-pipeline.bucket
  key    = "${var.platform}/config.yaml"
  source = var.config_file
}

resource "aws_s3_bucket_object" "trusted_data_file" {
  bucket = aws_s3_bucket.secure-pipeline.bucket
  key    = "${var.platform}/trusted_data.json"
  source = var.trusted_data_file
}

resource "aws_s3_bucket_object" "policies" {
  bucket = aws_s3_bucket.secure-pipeline.bucket
  key    = "${var.platform}/policies/${each.value}"
  source = "${var.policies_dir}/${each.value}"
  for_each = fileset(var.policies_dir, "*.rego")
}
