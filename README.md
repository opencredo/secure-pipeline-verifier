# Secure CI/CD Pipeline with policy-as-code
### This service makes use of GitHub REST APIs to collect information related to a GitHub Repository   
###
#### GitHub and Environment Configuration:
- User needs to have Admin access to the repository for some APIs to return correct response.
- GitHub Personal Access Token: You need to define a personal access token on GitHub with repository scope enabled for the APIs to return information of your repository.
- Create an environment variable called *GITHUB_TOKEN* with the GitHub Personal Access Token you generated as its value.

###
#### Service Configuration
The service has 2 arguments:  
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
    "github_repo": "opencredo/spring-cloud-stream",
    "pipeline_type": "travis",
    "trusted_users": [
      "afaedda"
    ]
  }
}

````

- The second parameter is a date with format *"YYYY-MM-ddTHH:mm:ss.SSSZ"* since when you want to check activity on your repository.

###
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