# Control-1 test
package github.user.cicd.auth

config = {
            "github_repo": "oc-org/my-cool-app",
            "pipeline_type": ".travis.yaml",
            "trusted_users": [
              "jsmith",
              "ajohnson"
            ]
         }

config_wrong_repo = {
            "github_repo": "oc-org/my-cool-application",
            "pipeline_type": ".travis.yaml",
            "trusted_users": [
              "jsmith",
              "ajohnson"
            ]
         }

test_authorized_cicd_change {
    safe_commit_input := {
                          "GitHubRepo": "oc-org/my-cool-app",
                          "CommitUrl": "https://github.com/oc-org/my-cool-app/commit/43565gt5464hbtr7665hty",
                          "Date": "2019-01-14T16:47:46Z",
                          "AuthorName": "John Smith",
                          "AuthorUsername": "jsmith",
                          "AuthorEmail": "jsmith@gmail.com",
                          "VerifiedSignature": true,
                          "VerificationReason": "valid"
                        }

    expected := "Control 1: INFO - Commit to CI/CD pilepine on repo [oc-org/my-cool-app] from user [jsmith] is authorized."

    is_authorized[expected] with input as safe_commit_input with data.config as config
}

test_unauthorized_cicd_change {

    unsafe_commit_input := {
                       "GitHubRepo": "oc-org/my-cool-app",
                       "CommitUrl": "https://github.com/oc-org/my-cool-app/commit/fvrer565eb564uh54",
                       "Date": "2021-06-20T01:47:46Z",
                       "AuthorName": "James K",
                       "AuthorUsername": "jk1234",
                       "AuthorEmail": "jk1234@gmail.com",
                       "VerifiedSignature": false,
                       "VerificationReason": "unsigned"
                     }

    expected := "Control 1: WARNING - User [jk1234] was not authorized to make changes to CI/CD on project repo [oc-org/my-cool-app]. Check commit details: https://github.com/oc-org/my-cool-app/commit/fvrer565eb564uh54"

    is_authorized[expected] with input as unsafe_commit_input with data.config as config
}

test_wrong_repo_config {

    commit_input := {
                        "GitHubRepo": "oc-org/my-cool-app",
                        "CommitUrl": "https://github.com/oc-org/my-cool-app/commit/43565gt5464hbtr7665hty",
                        "Date": "2019-01-14T16:47:46Z",
                        "AuthorName": "John Smith",
                        "AuthorUsername": "jsmith",
                        "AuthorEmail": "jsmith@gmail.com",
                        "VerifiedSignature": true,
                        "VerificationReason": "valid"
                     }

    expected := "Control 1: ERROR - Input repo [oc-org/my-cool-app] differs from config repo [oc-org/my-cool-application]. Please check configuration data"

    is_authorized[expected] with input as commit_input with data.config as config_wrong_repo
}