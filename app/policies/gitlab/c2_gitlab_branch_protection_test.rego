#Control-2 test
package gitlab.branch.protection

test_branch_protected {
    protected_branch_input := {
                           "Repo": "oc-org/my-app-repo",
                           "BranchName": "master",
                           "SignatureProtected": true,
                         }

    expected := "INFO - The branch [master] of repository [oc-org/my-app-repo] is protected with signed commits as expected."

    is_protected[expected] with input as protected_branch_input
}

test_branch_unprotected {
    unprotected_branch_input := {
                                  "Repo": "oc-org/my-app-repo",
                                  "BranchName": "develop",
                                  "SignatureProtected": false,
                                }

    expected := "WARNING - The branch [develop] of repository [oc-org/my-app-repo] is not protected with signed commits as expected. Please consider protecting it."

    is_protected[expected] with input as unprotected_branch_input
}