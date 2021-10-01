package github.user.cicd.auth

import data.config

default control = "Control 1"
default message = ""

is_authorized[message] {
    message := verify(input, config)
}

verify(commitDetails, configData) = response {
    commitDetails.Repo == configData.repo
    not user_authorized(commitDetails.AuthorUsername, configData.trusted_users)
    response := sprintf("%v: WARNING - User [%v] was not authorized to make changes to CI/CD on project repo [%v]. Check commit details: %v",
        [control, commitDetails.AuthorUsername, commitDetails.Repo, commitDetails.CommitUrl]
    )
}

verify(commitDetails, configData) = response {
    commitDetails.Repo == configData.repo
    user_authorized(commitDetails.AuthorUsername, configData.trusted_users)
    response := sprintf("%v: INFO - Commit to CI/CD pilepine on repo [%v] from user [%v] is authorized.",
        [control, commitDetails.Repo, commitDetails.AuthorUsername]
    )
}

verify(commitDetails, configData) = response {
    commitDetails.Repo != configData.repo
    response := sprintf("%v: ERROR - Input repo [%v] differs from config repo [%v]. Please check configuration data",
        [control, commitDetails.Repo, configData.repo]
    )
}

user_authorized(authorUsername, trustedUsers) {
    authorUsername == trustedUsers[_]
}
