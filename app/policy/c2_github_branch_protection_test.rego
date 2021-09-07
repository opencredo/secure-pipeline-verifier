#Control-2 test
package github.branch.protection

test_branch_protected {
    is_unprotected with input as {
                           "GitHubRepo": "oc-org/my-app-repo",
                           "BranchName": "master",
                           "SignatureProtected": true,
                           "Error": ""
                         }
}

test_branch_unprotected {
    unprotected_branch_input := {
                                  "GitHubRepo": "oc-org/my-app-repo",
                                  "BranchName": "develop",
                                  "SignatureProtected": false,
                                  "Error": "GET https://api.github.com/repos/oc-org/my-app-repo/branches/develop/protection/required_signatures: 404 Branch not protected []"
                                }

    expected := "WARNING - The branch [develop] of repository [oc-org/my-app-repo] is not protected with signed commits as expected. Please consider protecting it."

    is_unprotected[expected] with input as unprotected_branch_input
}