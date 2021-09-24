#Control-2 test
package gitlab.repo.protection

test_repo_protected {
    protected_repo_input := {
                           "Repo": "oc-org/my-app-repo",
                           "SignatureProtected": true,
                         }

    expected := "INFO - The repository [oc-org/my-app-repo] is protected with signed commits as expected."

    is_protected[expected] with input as protected_repo_input
}

test_repo_unprotected {
    unprotected_repo_input := {
                                  "Repo": "oc-org/my-app-repo",
                                  "SignatureProtected": false,
                                }

    expected := "WARNING - The repository [oc-org/my-app-repo] is not protected with signed commits as expected. Please consider protecting it."

    is_protected[expected] with input as unprotected_repo_input
}