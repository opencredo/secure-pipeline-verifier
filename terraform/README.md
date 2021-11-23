# Secure Pipeline: Terraform AWS Provisioner

## Overview
Use this Terraform config to provision AWS resources and run Secure Pipeline in AWS Lambda 

## How to run:

1. Define the following parameters. For example in `terraform.tfvars`:
```terraform
bucket = "bucket-name"
lambda_function_name="<Optional: name of the lambda function"
lambda_zip_file="<path to the zip file>"
region="<Default: eu-west-2. Region for the AWS resources>"
slack_token="<Using this value as an environment variable is recommended (i.e. TF_VAR_slack_token=<token>)"
lambda_timeout="Default: 3. Timeout (in seconds) for the lambda function."
repo_list = [{
    path       = "<path to a dir containing config and policy files for a repository>",
    repo_token = "[REDACTED]"
    event_schedule_rate = "Default: rate(12). This field is optional."
  },
  {
    path       = "<path to dir>",
    repo_token = "[REDACTED]"
  }]
```
2. Generate plan: `terraform plan`
3. Apply `terraform apply`