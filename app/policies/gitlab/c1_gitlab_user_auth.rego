package gitlab.user.cicd.auth

import data.config

default message = ""
default control = "Control 1"

is_authorized[message] {
    message := verify(input, config)
}

verify(commitDetails, configData) = response {
    commitDetails.Repo == configData.repo
    not user_authorized(commitDetails.AuthorName, configData.trusted_users)
    response := sprintf("%v: WARNING - User [%v] was not authorized to make changes to CI/CD on project repo [%v]. Check commit details: %v",
        [control, commitDetails.AuthorName, commitDetails.Repo, commitDetails.CommitUrl]
    )
}

verify(commitDetails, configData) = response {
    commitDetails.Repo == configData.repo
    user_authorized(commitDetails.AuthorName, configData.trusted_users)
    response := sprintf("%v: INFO - Commit to CI/CD pilepine on repo [%v] from user [%v] is authorized.",
        [control, commitDetails.Repo, commitDetails.AuthorName]
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
