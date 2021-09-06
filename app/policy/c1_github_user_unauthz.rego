package github.user.cicd.auth

import data.config

default message = "Allowed"

is_unauthorized[message] {
    input.GitHubRepo == config.github_repo
    input.AuthorUsername != config.trusted_users[_]
    message := sprintf("WARNING - User [%v] was not authorized to make changes to CI/CD on project repo [%v]. Check commit details: %v", [input.AuthorUsername, input.GitHubRepo, input.CommitUrl])
}