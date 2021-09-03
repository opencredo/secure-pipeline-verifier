package github.branch.protection

test_branch_protected {
    allow with input as {
                           "GitHubRepo": "oc-org/my-app-repo",
                           "BranchName": "master",
                           "SignatureProtected": true,
                           "Error": ""
                         }
}

test_branch_unprotected {
    not allow with input as {
                              "GitHubRepo": "",
                              "BranchName": "develop",
                              "SignatureProtected": false,
                              "Error": "GET https://api.github.com/repos/oc-org/my-app-repo/branches/develop/protection/required_signatures: 404 Branch not protected []"
                            }
}