## Configuration
To be able to run the application, we need to create two configuration files:

- *config.yaml* contains general configuration to run the application
- *trusted-data.yaml* contains private information to be used together with the policy controls

### config.yaml
This file is composed of a few blocks of configuration: 

```yaml
project:
  platform: github/gitlab
  owner: org-name
  repo: repo-name

repo-info-checks:
  ci-cd-path: path/to/ci-cd/config        # e.g. .travis.yaml, Jenkinsfile, .github/workflows, .gitlab-ci.yaml
  policies:
    - control: c1
      enabled: true/false
      path: "<path-to-policies>/cicd-auth.rego"
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
    - master
    - develop
    - other-branch

notifications:
  slack:
    enabled: true/false
    level: INFO/WARNING/ERROR
    notification-channel: <my-slack-notification-channel>
```

The *project* block defines the location of the repository that we want to run the policy checks on.
Here we have:

- **platform**: the platform where the repository is stored (github, gitlab)
- **owner**: the organization or profile name 
- **repo**: the actual repository or project name

The *repo-info-checks* block contains some details about the repository and the policy controls we want to run.
Under this block we have these configuration details: 

- **ci-cd-path**: the path to the CI/CD pipeline configuration file(s). 
- **policies**: this block contains information on the policies to run:  
  - **control**: this is just a static naming convention used internally in the application
  - **enabled**: enable/disable the run for this policy
  - **path**: the path to the policy definition file
- **protected-branches**: This defines the list of the branches that we expect to be protected with signed commits on our repository

The *notifications* block contains configuration for application notifications.  
Here we have the **slack** block with configuration regarding **slack** notifications:

- **enabled**: enable/disable slack notifications
- **level**: the level of notification you want to receive. 
  - **INFO** notifies you that your repo meets the requirement for the given policy
  - **WARNING** notifies you that your repo doesn't meet the requirement for the given policy
  - **ERROR** notifies that there's been an error while checking the policy
- **notification-channel**: this is the name of your slack channel where you want to receive the notifications

### trusted-data.yaml

This file contains the list of users that have authorization to make changes to the CI/CD pipeline configuration on your repo

```yaml
config:
  repo: some-org/some-repo
  pipeline_type: travis
  trusted_users:
    - username-1
    - username-2
    - username-3
```

Very simply, the *config* block contains the following fields: 

- **repo**: The full name of your repo
- **pipeline_type**: the type of CI/CD pipeline tool you're using
- **trusted_users**: the list of users with authorization to the CI/CD pipeline

### Environment Variables

In case you're running this application on the CLI, then you'll also need to set the following Environment Variables: 

- **GITHUB_TOKEN**: If your repository is on GitHub, you need to define this variable so that the application is able to make use of the APIs
- **GITLAB_TOKEN**: If your repository is on GitLab, you need to define this variable so that the application is able to make use of the APIs
- **SLACK_TOKEN**: If you enable the notifications on Slack, you need to set this token to be able to connect to your Slack channel.

If you're running this application as a Lambda function on AWS, you'll need to set these on the Terraform vars.
You can find a guide on how to do so [here](../terraform/README.md).