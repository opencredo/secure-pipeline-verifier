# Secure CI/CD Pipeline with policy-as-code
### This service makes use of git repositories REST APIs to collect information related to a repository

### Supported git repositories:

### Github
#### Environment Configuration:
  - User needs to have Admin access to the repository for some APIs to return correct response.
  - GitHub Personal Access Token: You need to define a personal access token on GitHub with repository scope enabled for the APIs to return information of your repository.
  - Create an environment variable called *GITHUB_TOKEN* with the GitHub Personal Access Token you generated as its value.

### Gitlab
#### Environment Configuration:
- User needs to have Admin access to the repository for some APIs to return correct response.
- GitLab Personal Access Token: You need to define a personal access token on GitLab with **read_api** and **read_repository** scopes.
- Create an environment variable called *GITLAB_TOKEN* with the Gitlab Personal Access Token you generated as its value.


### Service Configuration
#### The service can run as a local standalone service from command-line or as AWS Lambda function

#### Standalone command-line configuration
The service has 2 arguments:  
- First argument is the path to a yaml file with the following structure:  

````
project:
    platform: <github/gitlab>
    owner: <your repo organisation name>
    repo: <your repo name>

repo-info-checks:
    ci-cd-path: <path to your CI/CD pipeline>
    policies:
        - control: c1
          enabled: true/false
          path: "<path-to-policies>/auth.rego"
        - control: c2
          enabled: true/false
          path: "<path-to-policies>/signed-commits.rego"
        - control: c3
          enabled: true/false
          path: "<path-to-policies>/auth-key-expiry.rego"
        - control: c4
          enabled: true/false
          path: "<path-to-policies>/auth-key-read-only.rego"
    protected-branches:
        - branch-1
        - branch-2
slack:
    enabled: true/false
    notification-channel: "<channel-name>"
````

Here's an example of this configuration:
````
project:
    platform: github
    owner: opencredo
    repo: spring-cloud-stream

repo-info-checks:
    ci-cd-path: .travis.yaml
    policies:
        - control: c1
          enabled: true
          path: "/policies/ci-auth.rego"
        - control: c2
          enabled: false
          path: "/policies/protected-commits.rego"
        - control: c3
          enabled: false
          path: "/policies/key-auto-expire.rego"
        - control: c4
          enabled: true
          path: "/policies/key-read-only.rego"
    protected-branches:
        - master
        - develop

slack:
    enabled: true
    notification-channel: "secure-pipeline-notifications"

````

The trusted-data-file is a json file acting as a source of truth about your repository. 
Here's an example: 

````
{
  "config": {
    "repo": "opencredo/spring-cloud-stream",
    "pipeline_type": "travis",
    "trusted_users": [
      "afaedda"
    ]
  }
}

````

- The second parameter is a date with format *"YYYY-MM-ddTHH:mm:ss.SSSZ"* since when you want to check activity on your repository.

- (Optional) The third parameter is the name of the branch you want to audit. 

#### AWS Lambda configuration

To deploy the service as a Lambda function and the required infrastructure for execution, check the `terraform/README.md` file.

The AWS Lambda function expects a JSON input event with the following structure:

````
{
    "region": "<aws-region>"
    "bucket": "<s3-bucket-containing-the-configuration>"
    "configPath": "<bucket-path-to-config>"
    "branch": "<(optional) branch's name>"
}
````

This JSON input is created as part of the Terraform infrastructure configuration

#### Testing

To run OPA rego tests from *app/policy* directory you need to: 
- Install OPA on your local machine  
````
    $ brew install opa 
````
 - Run the tests with the command
````
    opa test . -v
````