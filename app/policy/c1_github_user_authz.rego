# Control-1

package github.cicd.user.authz

import data.config

default allow = false

allow {
    input.GitHubRepo == config.github_repo
    input.AuthorUsername == config.trusted_users[_]
}
