#Control-2 test
package github.branch.protection

test_branch_protected {
    protected_branch_input := {
                           "GitHubRepo": "oc-org/my-app-repo",
                           "BranchName": "master",
                           "SignatureProtected": true,
                           "Error": ""
                         }

    expected := "Control 2: INFO - The branch [master] of repository [oc-org/my-app-repo] is protected with signed commits as expected."

    is_protected[expected] with input as protected_branch_input
}

test_branch_unprotected {
    unprotected_branch_input := {
                                  "GitHubRepo": "oc-org/my-app-repo",
                                  "BranchName": "develop",
                                  "SignatureProtected": false,
                                  "Error": "GET https://api.github.com/repos/oc-org/my-app-repo/branches/develop/protection/required_signatures: 404 Branch not protected []"
                                }

    expected := "Control 2: WARNING - The branch [develop] of repository [oc-org/my-app-repo] is not protected with signed commits as expected. Please consider protecting it."

    is_protected[expected] with input as unprotected_branch_input
}

test_not_permitted {
    user_not_permitted_input := {
                                  "GitHubRepo": "oc-org/my-app-repo",
                                  "BranchName": "develop",
                                  "SignatureProtected": false,
                                  "Error": "GET https://api.github.com/repos/oc-org/my-app-repo/branches/develop/protection/required_signatures: 404 Not Found []"
                                }

    expected := "Control 2: ERROR - The user has not Admin permissions on repository [oc-org/my-app-repo] to perform this check. Please consider updating permissions."

    is_protected[expected] with input as user_not_permitted_input
}