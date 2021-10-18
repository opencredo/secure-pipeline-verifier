#Control-2 test
package signature.protection

test_repo_protected {
    protected_repo_input := {
                           "Repo": "oc-org/my-app-repo",
                           "SignatureProtected": true,
                         }

    expected := {
        "control": "Control 2",
        "level": "INFO",
        "msg": "The repository [oc-org/my-app-repo] is protected with signed commits as expected."
    }

    decision := is_protected with input as protected_repo_input

    decision.control == "Control 2"
    decision.level == "INFO"
    decision.msg == "The repository [oc-org/my-app-repo] is protected with signed commits as expected."
}

test_repo_unprotected {
    unprotected_repo_input := {
                                  "Repo": "oc-org/my-app-repo",
                                  "SignatureProtected": false,
                                }

    expected := {
        "control": "Control 2",
        "level": "WARNING",
        "msg": "The repository [oc-org/my-app-repo] is not protected with signed commits as expected. Please consider protecting it."
    }

    decision := is_protected with input as unprotected_repo_input

    decision.control == "Control 2"
    decision.level == "WARNING"
    decision.msg == "The repository [oc-org/my-app-repo] is not protected with signed commits as expected. Please consider protecting it."
}