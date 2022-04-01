# Secure CI/CD Pipeline with policy-as-code
### This service makes use of git repositories REST APIs to collect information related to a repository

### Supported git repositories:

### Github
#### Environment Configuration:
  - User needs to have Admin access to the repository for some APIs to return correct response.
  - GitHub Personal Access Token: You need to define a personal access token on GitHub with repository scope enabled for the APIs to return information of your repository.
  - Create an environment variable called *REPO_TOKEN* with the GitHub Personal Access Token you generated as its value.

### Gitlab
#### Environment Configuration:
- User needs to have Admin access to the repository for some APIs to return correct response.
- GitLab Personal Access Token: You need to define a personal access token on GitLab with **read_api** and **read_repository** scopes.
- Create an environment variable called *REPO_TOKEN* with the Gitlab Personal Access Token you generated as its value.


### Service Configuration
#### The service can run as a local standalone service from command-line or as AWS Lambda function

#### Standalone command-line configuration
The service has several arguments:  
- First argument is the path to a yaml file with the following structure:  

````
project:
  platform: github/gitlab
  owner: org-name
  repo: repo-name

repo-info:
  ci-cd-path: <path to your CI/CD pipeline>
  protected-branches:
    - master
    - develop
    - other-branch

policies:
  - control: c1
    enabled: true/false
    path: "<path to policies>/auth.rego"
  - control: c2
    enabled: true/false
    path: "<path to policies>/signed-commits.rego"
  - control: c3
    enabled: true/false
    path: "<path to policies>/auth-key-expiry.rego"
  - control: c4
    enabled: true/false
    path: "<path to policies>/auth-key-read-only.rego"

notifications:
  slack:
    enabled: true/false
    level: INFO/WARNING/ERROR
    notification-channel: "<channel-name>"
````

Here's an example of this configuration:
````
project:
    platform: github
    owner: opencredo
    repo: spring-cloud-stream

repo-info:
  ci-cd-path: ".travis.yml"
  protected-branches:
    - master
    - develop
    - other-branch

policies:
  - control: c1
    enabled: true/false
    path: "<path to policies>/auth.rego"
  - control: c2
    enabled: true/false
    path: "<path to policies>/signed-commits.rego"
  - control: c3
    enabled: true/false
    path: "<path to policies>/auth-key-expiry.rego"
  - control: c4
    enabled: true/false
    path: "<path to policies>/auth-key-read-only.rego"

notifications:
  slack:
    enabled: true
    level: INFO
    notification-channel: "secure-pipeline-notifications"

````

The trusted-data-file is a yaml file acting as a source of truth about your repository. 
Here's an example: 

````
config:
  repo: opencredo/spring-cloud-stream
  pipeline_type: travis
  trusted_users:
    - afaedda
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