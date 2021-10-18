package user.cicd.auth

import data.config

default control = "Control 1"

is_authorized = decision {
    decision := verify(input, config)
}

verify(commitDetails, configData) = decision {
    commitDetails.Repo == configData.repo
    not user_authorized(commitDetails.AuthorName, configData.trusted_users)
    response := sprintf("User [%v] was not authorized to make changes to CI/CD on project repo [%v]. Check commit details: %v",
        [commitDetails.AuthorName, commitDetails.Repo, commitDetails.CommitUrl]
    )

    decision := {"control": control, "level": "WARNING", "msg": response}
}

verify(commitDetails, configData) = decision {
    commitDetails.Repo == configData.repo
    user_authorized(commitDetails.AuthorName, configData.trusted_users)
    response := sprintf("Commit to CI/CD pilepine on repo [%v] from user [%v] is authorized.",
        [commitDetails.Repo, commitDetails.AuthorName]
    )

    decision := {"control": control, "level": "INFO", "msg": response}
}

verify(commitDetails, configData) = decision {
    commitDetails.Repo != configData.repo
    response := sprintf("Input repo [%v] differs from config repo [%v]. Please check configuration data",
        [commitDetails.Repo, configData.repo]
    )

    decision := {"control": control, "level": "ERROR", "msg": response}
}

user_authorized(authorUsername, trustedUsers) {
    authorUsername == trustedUsers[_]
}
