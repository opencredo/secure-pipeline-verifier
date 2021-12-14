# Secure Pipelines
The *Secure Pipelines* application runs policy checks on a repository hosted in GitHub or GitLab platform.

## Policies
Here's the list of policies checks that the application allows you to run on your repository: 

### Control 1: Restrict administrative access to CI/CD tools
It's important to ensure that only authorized persons can make administrative changes to the CI/CD system. 
If an unauthorized person is able to gain access, changes to pipeline definitions enable the subversion of many of the remaining controls in this document.

### Control 2: Only accept commits signed with a developer GPG key
Unsigned code commits are challenging, if not impossible, to trace and pose a risk to the integrity of the code base. 
Requiring commits to be signed with a developer GPG key helps to ensure non-repudiation of commits and increases the burden on the attacker seeking to insert malicious code.

### Control 3: Automation access keys expire automatically
Ensuring that access keys used by automation expire periodically creates a shorter window of attack when keys are compromised.

### Control 4: Reduce automation access to read-only
CI systems should have *read only access* to source code repositories following the principle of least privilege access.

## Usage
In order to use this application, you must have an active account on GitHub/GilLab and admin access to the repository you want to run the policies on.
You also need to generate a Personal Access Token on [GitHub](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) 
or [GitLab](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html) with repo access for the APIs to work.

In this guide we're going to take a look at how to [configure](config.md) it, [deploy](deploy.md) it and what we can [expect](notifications.md) from it.