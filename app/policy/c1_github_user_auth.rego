package github.user.cicd.auth

import data.config

default message = ""

is_authorized[message] {
    message := verify(input, config)
}

verify(commitDetails, configData) = response {
    commitDetails.GitHubRepo == configData.github_repo
    not user_authorized(commitDetails.AuthorUsername, configData.trusted_users)
    response := sprintf("WARNING - User [%v] was not authorized to make changes to CI/CD on project repo [%v]. Check commit details: %v", [commitDetails.AuthorUsername, commitDetails.GitHubRepo, commitDetails.CommitUrl])
}

verify(commitDetails, configData) = response {
    commitDetails.GitHubRepo == configData.github_repo
    user_authorized(commitDetails.AuthorUsername, configData.trusted_users)
    response := sprintf("INFO - Commit to CI/CD pilepine on repo [%v] from user [%v] is authorized.", [commitDetails.GitHubRepo, commitDetails.AuthorUsername])
}

verify(commitDetails, configData) = response {
    commitDetails.GitHubRepo != configData.github_repo
    response := sprintf("ERROR - Input repo [%v] differs from config repo [%v]. Please check configuration data", [commitDetails.GitHubRepo, configData.github_repo])
}

user_authorized(authorUsername, trustedUsers) {
    authorUsername == trustedUsers[_]
}
