package github.user.cicd.auth

config = {
              "github_repo": "oc-org/my-cool-app",
              "pipeline_type": ".travis.yaml",
              "trusted_users": [
                "jsmith",
                "ajohnson"
              ]
         }

test_authorized_cicd_change {
    is_unauthorized with input as {
                          "GitHubRepo": "oc-org/my-cool-app",
                          "CommitUrl": "https://github.com/oc-org/my-cool-app/commit/43565gt5464hbtr7665hty",
                          "Date": "2019-01-14T16:47:46Z",
                          "AuthorName": "John Smith",
                          "AuthorUsername": "jsmith",
                          "AuthorEmail": "jsmith@gmail.com",
                          "VerifiedSignature": true,
                          "VerificationReason": "valid"
                        }
    with data.config as config
}

test_is_unauthorized_cicd_change {

    unsafe_commit := {
                       "GitHubRepo": "oc-org/my-cool-app",
                       "CommitUrl": "https://github.com/oc-org/my-cool-app/commit/fvrer565eb564uh54",
                       "Date": "2021-06-20T01:47:46Z",
                       "AuthorName": "James K",
                       "AuthorUsername": "jk1234",
                       "AuthorEmail": "jk1234@gmail.com",
                       "VerifiedSignature": false,
                       "VerificationReason": "unsigned"
                     }

    expected := "User jk1234 was not authorized to make changes to CI/CD on project repo oc-org/my-cool-app. Check commit details: https://github.com/oc-org/my-cool-app/commit/fvrer565eb564uh54"

    is_unauthorized[expected] with input as unsafe_commit with data.config as config
}