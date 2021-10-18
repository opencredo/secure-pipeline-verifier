# Control-1 test
package user.cicd.auth

config = {
            "repo": "oc-org/my-cool-app",
            "pipeline_type": ".travis.yaml",
            "trusted_users": [
              "jsmith",
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
                          "CommitUrl": "https://github.com/oc-org/my-cool-app/commit/43565gt5464hbtr7665hty",
                          "Date": "2019-01-14T16:47:46Z",
                          "AuthorName": "John Smith",
                          "AuthorUsername": "jsmith",
                          "AuthorEmail": "jsmith@gmail.com",
                          "VerifiedSignature": true,
                          "VerificationReason": "valid"
                        }

    expected := {
        "control": "Control 1",
        "level": "INFO",
        "msg": "Commit to CI/CD pipeline on repo [oc-org/my-cool-app] from user [jsmith] is authorized."
    }

    decision := is_authorized with input as safe_commit_input with data.config as config

    decision.control == "Control 1"
    decision.level == "INFO"
    decision.msg == "Commit to CI/CD pipeline on repo [oc-org/my-cool-app] from user [jsmith] is authorized."
}

test_unauthorized_cicd_change {

    unsafe_commit_input := {
                       "Repo": "oc-org/my-cool-app",
                       "CommitUrl": "https://github.com/oc-org/my-cool-app/commit/fvrer565eb564uh54",
                       "Date": "2021-06-20T01:47:46Z",
                       "AuthorName": "James K",
                       "AuthorUsername": "jk1234",
                       "AuthorEmail": "jk1234@gmail.com",
                       "VerifiedSignature": false,
                       "VerificationReason": "unsigned"
                     }

    expected := {
        "control": "Control 1",
        "level": "WARNING",
        "msg": "User [jk1234] was not authorized to make changes to CI/CD on project repo [oc-org/my-cool-app]. Check commit details: https://github.com/oc-org/my-cool-app/commit/fvrer565eb564uh54"
    }

    decision := is_authorized with input as unsafe_commit_input with data.config as config

    decision.control == "Control 1"
    decision.level == "WARNING"
    decision.msg == "User [jk1234] was not authorized to make changes to CI/CD on project repo [oc-org/my-cool-app]. Check commit details: https://github.com/oc-org/my-cool-app/commit/fvrer565eb564uh54"
}

test_wrong_repo_config {

    commit_input := {
                        "Repo": "oc-org/my-cool-app",
                        "CommitUrl": "https://github.com/oc-org/my-cool-app/commit/43565gt5464hbtr7665hty",
                        "Date": "2019-01-14T16:47:46Z",
                        "AuthorName": "John Smith",
                        "AuthorUsername": "jsmith",
                        "AuthorEmail": "jsmith@gmail.com",
                        "VerifiedSignature": true,
                        "VerificationReason": "valid"
                     }

    expected := {
        "control": "Control 1",
        "level": "ERROR",
        "msg": "Input repo [oc-org/my-cool-app] differs from config repo [oc-org/my-cool-application]. Please check configuration data"
    }

    decision := is_authorized with input as commit_input with data.config as config_wrong_repo
    decision.control == "Control 1"
    decision.level == "ERROR"
    decision.msg == "Input repo [oc-org/my-cool-app] differs from config repo [oc-org/my-cool-application]. Please check configuration data"
}