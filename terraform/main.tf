provider "aws" {
    region = "eu-west-2"
  }

terraform {
  backend "s3" {
    bucket = "alberto-secure-pipeline-tf-state"
    key    = "af-terraform-state"
    region = "eu-west-2"
  }
}

resource "aws_s3_bucket" "config-bucket" {
  bucket = "secure-pipelines-config"
  acl = "private"

  tags = {
    Name = "Root for secure-pipelines repos config"
  }
}

resource "aws_s3_bucket_object" "config" {
  for_each = fileset("../../config/", "*")
  bucket = aws_s3_bucket.config-bucket.id
  key = "opencredo/spring-cloud-stream/${each.value}"
  source = "../../config/${each.value}"
  etag = filemd5("../../config/${each.value}")
}