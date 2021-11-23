# Module: Repository

## Overview

This module is used for provisioning the AWS infrastructure with the resources needed for auditing a repository 
by the Secure Pipeline application.

## Usage

```terraform
# main.tf

(...)
        
module "repositories" {
  source = "./modules/repository"
  for_each = { for repo in var.repo_list : repo.path => repo }
  source_dir  = each.key
  bucket      = aws_s3_bucket.secure_pipeline.bucket
  lambda_arn  = aws_lambda_function.check_policies.arn
  lambda_name = aws_lambda_function.check_policies.function_name
  last_run    = timestamp()
  parameter_prefix = var.parameter_prefix
  repo_token  = each.value.repo_token
  region = var.region
}
```