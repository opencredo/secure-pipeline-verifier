## Deployment

The application is organized in a way so that it could be run as a CLI executable or a Lambda function.  
This organization is reflected on its folder structure where we have two entry points under the *cmd* folder.  

- The AWS Lambda entrypoint is *cmd/aws/main.go*  
- The CLI entrypoint is *cmd/cli/main.go*

### CLI 

To run the application from the CLI, you need to provide two arguments:
- The path to your configuration files (config.yaml, trusted-data.yaml)
- The date-time since when you want the policy controls to check your repo

You can find a guide for the configurations [here](config.md). 

The following is an example on how to run it from the CLI: 

```shell
$ cd cmd/cli
$ go run main.go "/path/to/config/" "2020-01-01T09:00:00.000Z" branch_name
```

### Terraform

Terraform is used in this project for provisioning an AWS infrastructure for cloud deployment of the application.
This allows the Secure Pipeline Verifier to be run on a schedule, and it can be invoked via an API call.
Additionally, with this deployment the users have access to ChatOps which is a way to run the application via a Slack command.
In  [main.tf](../terraform/main.tf) you can check which resources will be created by Terraform.
