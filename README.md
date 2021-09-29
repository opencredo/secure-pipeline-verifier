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


#### Service Configuration
The service has 3 arguments:  
- First argument is the path to a yaml file with the following structure:  

````
project:
    owner: <your repo organisation name>
    repo: <your repo name>

repo-info-checks:
    trusted-data-file: <path to your json trusted data>
    ci-cd-path: <path to your CI/CD pipeline>
    protected-branches:
        - branch-1
        - branch-2
````

Here's an example of this configuration:
````
project:
    owner: opencredo
    repo: spring-cloud-stream

repo-info-checks:
    trusted-data-file: prj-trusted-data.json
    ci-cd-path: .travis.yaml
    protected-branches:
        - master
        - develop
````

The trusted-data-file is a json file acting as a source of truth about your repository. 
Here's an example: 

````
{
  "config": {
    "repo": "org-name/repo-name",
    "pipeline_type": "travis",
    "trusted_users": [
      "some-name"
    ]
  }
}

````

- The second parameter is a date with format *"YYYY-MM-ddTHH:mm:ss.SSSZ"* since when you want to check activity on your repository.
- The third parameter is a type of git repository (i.e. **github**).

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