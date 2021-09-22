package gitlab.user.cicd.auth

import data.config

default message = ""

is_authorized[message] {
    message := verify(input, config)
}

verify(commitDetails, configData) = response {
    commitDetails.Repo == configData.repo
    not user_authorized(commitDetails.AuthorName, configData.trusted_users)
    response := sprintf("WARNING - User [%v] was not authorized to make changes to CI/CD on project repo [%v]. Check commit details: %v",
        [commitDetails.AuthorName, commitDetails.Repo, commitDetails.CommitUrl]
    )
}

verify(commitDetails, configData) = response {
    commitDetails.Repo == configData.repo
    user_authorized(commitDetails.AuthorName, configData.trusted_users)
    response := sprintf("INFO - Commit to CI/CD pilepine on repo [%v] from user [%v] is authorized.",
        [commitDetails.Repo, commitDetails.AuthorName]
    )
}

verify(commitDetails, configData) = response {
    commitDetails.Repo != configData.repo
    response := sprintf("ERROR - Input repo [%v] differs from config repo [%v]. Please check configuration data",
        [commitDetails.Repo, configData.repo]
    )
}

user_authorized(authorUsername, trustedUsers) {
    authorUsername == trustedUsers[_]
}
