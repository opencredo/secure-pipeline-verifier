# Secure Pipeline: Terraform AWS Provisioner

## Overview
Use this Terraform config to provision AWS resources and run Secure Pipeline and (optionally) ChatOps in AWS Lambdas. 

## AWS Resources
Below are resources that will be created in AWS by Terraform.
* API Gateway: It creates two endpoints `/audit`. If ChatOps is enabled, then the `/chatops` endpoint will be created. 
* CloudWatch logs: Stores logs coming from the Lambda functions.
* IAM Roles: Roles for the Lambda functions.
* Lambda - ChatOps: It is designed to parse the Slack command and trigger the Lambda-Secure Pipeline Verifier function.  
* Lambda - Secure Pipeline Verifier: The application in a Lambda function.
* Parameter Store: Stores values such as: `last_run`, `repo_roken`, `slack_token` for the Lambda - SPV function.
* S3 Bucket: Used for storing config files and policy files to specific repositories.
so that Slack can send the command to that endpoint.

## Prerequisites:

### Access to AWS
To provision an AWS infrastructure Terraform needs to be able to authenticate with AWS.
There are various ways to do this. You pick the most suitable option for you in this documentation from the AWS provider:
[https://registry.terraform.io/providers/hashicorp/aws/latest/docs#authentication-and-configuration](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#authentication-and-configuration)

### AWS Lambda - Secure Pipeline Verifier

In order to be able to run the application as an AWS Lambda function, we first need to build its executable and compress it
to a zip file, by executing the following commands: 

```shell
$ make build-lambda
```

### (Optional) AWS Lambda - ChatOps

In order to be able to run ChatOps, we first need to build its executable and compress it
to a zip file, by executing the following commands: 
```shell
$ make build-lambda-chatops
```

## Run Terraform:

1. Make sure you're in the `terraform/` directory.
2. Define the following parameters. For example in `terraform.tfvars`:
```terraform
bucket = "bucket-name"
lambda_function_name="<Optional: name of the lambda function"
lambda_zip_file="<path to the zip file>"
lambda_chatops_zip_file = "Default: null | <path to the zip file>"
lambda_timeout="Default: 3. Timeout (in seconds) for the lambda function."
last_run    = "<Optional: Format: 'YYYY-MM-DD'T'hh:mm:ssZ'>"
region="<Default: eu-west-2. Region for the AWS resources>"
slack_token="<A token to authenticate with Slack>"
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
3. Generate plan: `terraform plan`
4. Apply `terraform apply`
5. The output should contain a URL to the API, the following endpoints can be accessed:
   * `/audit` - triggers the application via an API call.
   * `/chatops` - Parses the Slack command and triggers the application. 
      NOTE: This endpoint is created when the `lambda_chatops_zip_file` argument is provided to Terraform. 
   

### Configure ChatOps

To enable ChatOps in Slack, you need to [create an application](https://api.slack.com/apps/) via Slack settings. 
After creating the app, in the application's manifest you will be able to create a custom Slack command and point it to the new API endpoint.
