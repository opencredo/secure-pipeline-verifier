# Control-1 test
package gitlab.user.cicd.auth

config = {
            "repo": "oc-org/my-cool-app",
            "pipeline_type": ".travis.yaml",
            "trusted_users": [
              "John Smith",
              "ajohnson"
            ]
         }

config_wrong_repo = {
            "repo": "oc-org/my-cool-application",
            "pipeline_type": ".travis.yaml",
            "trusted_users": [
              "jsmith",
              "ajohnson"
            ]
         }

test_authorized_cicd_change {
    safe_commit_input := {
                          "Repo": "oc-org/my-cool-app",
                          "CommitUrl": "https://gitlab.example.com/thedude/gitlab-foss/-/commit/ed899a2f4b50b4370feeea94676502b42383c746",
                          "Date": "2019-01-14T16:47:46Z",
                          "AuthorName": "John Smith",
                          "AuthorEmail": "jsmith@gmail.com",
                          "VerifiedSignature": true,
                          "VerificationReason": "valid"
                        }

    expected := "INFO - Commit to CI/CD pilepine on repo [oc-org/my-cool-app] from user [John Smith] is authorized."

    is_authorized[expected] with input as safe_commit_input with data.config as config
}

test_unauthorized_cicd_change {

    unsafe_commit_input := {
                       "Repo": "oc-org/my-cool-app",
                       "CommitUrl": "https://gitlab.example.com/thedude/gitlab-foss/-/commit/ed899a2f4b50b4370feeea94676502b42383c746",
                       "Date": "2021-06-20T01:47:46Z",
                       "AuthorName": "James K",
                       "AuthorEmail": "jk1234@gmail.com",
                       "VerifiedSignature": false,
                       "VerificationReason": "unsigned"
                     }

    expected := "WARNING - User [James K] was not authorized to make changes to CI/CD on project repo [oc-org/my-cool-app]. Check commit details: https://gitlab.example.com/thedude/gitlab-foss/-/commit/ed899a2f4b50b4370feeea94676502b42383c746"

    is_authorized[expected] with input as unsafe_commit_input with data.config as config
}

test_wrong_repo_config {

    commit_input := {
                        "Repo": "oc-org/my-cool-app",
                        "CommitUrl": "https://gitlab.example.com/thedude/gitlab-foss/-/commit/ed899a2f4b50b4370feeea94676502b42383c746",
                        "Date": "2019-01-14T16:47:46Z",
                        "AuthorName": "John Smith",
                        "AuthorEmail": "jsmith@gmail.com",
                        "VerifiedSignature": true,
                        "VerificationReason": "valid"
                     }

    expected := "ERROR - Input repo [oc-org/my-cool-app] differs from config repo [oc-org/my-cool-application]. Please check configuration data"

    is_authorized[expected] with input as commit_input with data.config as config_wrong_repo
}