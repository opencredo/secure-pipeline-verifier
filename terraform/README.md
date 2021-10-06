# Secure Pipeline: Terraform AWS Provisioner

## Overview
Use this Terraform config to provision AWS resources and run Secure Pipeline in AWS Lambda 

## How to run:

1. Define the following parameters (For example in `terraform.tfvars`):
```terraform
bucket = "bucket-name"
config_file="<path to>/config.yaml"
platform="github"
policies_dir="<path to>/policies"
trusted_data_file="<path to>/trusted_data.json"
```
2. Generate plan: `terraform plan`
3. Apply `terraform apply`