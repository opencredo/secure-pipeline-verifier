## Deployment

The application is organized in a way so that it could be run as a CLI executable or a Lambda function.  
This organization is reflected on its folder structure where we have two entry points under the *cmd* folder.  

- The AWS Lambda entrypoint is *cmd/aws/main.go*  
- The CLI entrypoint is *cmd/cli/main.go*

### AWS Lambda (Secure Pipeline Verifier)

In order to be able to run the application as an AWS Lambda function, we first need to build its executable and compress it
to a zip file, by executing the following commands: 

```shell
$ make build-lambda
```

### AWS Lambda (Secure Pipeline Verifier)

In order to be able to run ChatOps, we first need to build its executable and compress it
to a zip file, by executing the following commands: 
```shell
$ make build-lambda-chatops
```

#### Terraform

After building and compressing the application executable, the next step is to configure and launch Terraform for infrastructure provisioning. 
You can find a guide on how to do so [here](../terraform/README.md).

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