# Secure Pipeline: Terraform AWS Provisioner

## Overview
Use this Terraform config to run Secure Pipeline as AWS Lambda 

## How to run:

1. Define the following parameters (For example in `terraform.tfvars`):
```terraform
platform="github"
policies_dir="<path to>/policies"
config_file="<path to>/config.yaml"
trusted_data_file="<path to>/trusted_data.json"
```
2. Generate plan: `terraform plan`
3. Apply `terraform apply`