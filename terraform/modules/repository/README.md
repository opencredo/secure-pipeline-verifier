# Module: Repository

## Overview

This module is used for provisioning the AWS infrastructure with the resources needed for auditing a repository 
by the Secure Pipeline application.

## Usage
- Provision the resources for multiple repositories:

```terraform
# terraform.tfvars

bucket          = "<bucket name>"
lambda_zip_file = "<path to>/function.zip"
repo_list = [{
    path       = "<path to dir>",
    repo_token = "[REDACTED]"
    event_schedule_rate = "[OPTIONAL - For example: rate(10)]"
  },
  {
    path       = "<path to dir>",
    repo_token = "[REDACTED]"
  }]
last_run    = "[OPTIONAL - Format: YYYY-MM-DD'T'hh:mm:ssZ]"
slack_token = "[REDACTED]"
```

```terraform
# main.tf

(...)

module "repositories" {
  source           = "./modules/repository"
  for_each         = { for repo in var.repo_list : repo.path => repo }
  source_dir       = each.key
  bucket           = aws_s3_bucket.secure_pipeline.bucket
  lambda_arn       = aws_lambda_function.check_policies.arn
  lambda_name      = aws_lambda_function.check_policies.function_name
  last_run         = coalesce(var.last_run, timestamp())
  parameter_prefix = var.parameter_prefix
  repo_token       = each.value.repo_token
  region           = var.region
}
```