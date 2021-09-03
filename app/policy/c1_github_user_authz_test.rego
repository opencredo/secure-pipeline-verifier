package github.cicd.user.authz

config = {
              "github_repo": "oc-org/my-cool-app",
              "pipeline_type": ".travis.yaml",
              "trusted_users": [
                "jsmith",
                "ajohnson"
              ]
         }



test_allowed_cicd_change {
    allow with input as {
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

test_not_allowed_cicd_change {
    not allow with input as {
                               "GitHubRepo": "oc-org/my-cool-app",
                               "CommitUrl": "https://github.com/oc-org/my-cool-app/commit/fvrer565eb564uh54",
                               "Date": "2021-06-20T01:47:46Z",
                               "AuthorName": "James K",
                               "AuthorUsername": "jk1234",
                               "AuthorEmail": "jk1234@gmail.com",
                               "VerifiedSignature": false,
                               "VerificationReason": "unsigned"
                             }
}