#Control-2 test
package signature.protection

test_branch_protected {
    protected_branch_input := {
                           "Repo": "oc-org/my-app-repo",
                           "BranchName": "master",
                           "SignatureProtected": true,
                           "Error": ""
                         }

    expected := {
        "control": "Control 2",
        "level": "INFO",
        "msg": "The branch [master] of repository [oc-org/my-app-repo] is protected with signed commits as expected."
    }

    decision := is_protected with input as protected_branch_input

    decision.control == "Control 2"
    decision.level == "INFO"
    decision.msg == "The branch [master] of repository [oc-org/my-app-repo] is protected with signed commits as expected."
}

test_branch_unprotected {
    unprotected_branch_input := {
                                  "Repo": "oc-org/my-app-repo",
                                  "BranchName": "develop",
                                  "SignatureProtected": false,
                                  "Error": "Branch not protected"
                                }

    expected := {
        "control": "Control 2",
        "level": "WARNING",
        "msg": "The branch [develop] of repository [oc-org/my-app-repo] is not protected with signed commits as expected. Please consider protecting it."
    }

    decision := is_protected with input as unprotected_branch_input

    decision.control == "Control 2"
    decision.level == "WARNING"
    decision.msg == "The branch [develop] of repository [oc-org/my-app-repo] is not protected with signed commits as expected. Please consider protecting it."
}

test_user_not_permitted {
    user_not_permitted_input := {
                                  "Repo": "oc-org/my-app-repo",
                                  "BranchName": "develop",
                                  "SignatureProtected": false,
                                  "Error": "Not Found"
                                }

    expected := {
        "control": "Control 2",
        "level": "ERROR",
        "msg": "The user has not Admin permissions on repository [oc-org/my-app-repo] to perform this check. Please consider updating permissions."
    }

    decision := is_protected with input as user_not_permitted_input

    decision.control == "Control 2"
    decision.level == "ERROR"
    decision.msg == "The user has not Admin permissions on repository [oc-org/my-app-repo] to perform this check. Please consider updating permissions."
}

test_not_existing_branch {
    not_existing_branch_input := {
                                  "Repo": "oc-org/my-app-repo",
                                  "BranchName": "not-existing-branch",
                                  "SignatureProtected": false,
                                  "Error": "Branch not found"
                                }

    expected := {
        "control": "Control 2",
        "level": "ERROR",
        "msg": "The branch [not-existing-branch] was not found in the repository [oc-org/my-app-repo]. Please check configuration."
    }

    decision := is_protected with input as not_existing_branch_input

    decision.control == "Control 2"
    decision.level == "ERROR"
    decision.msg == "The branch [not-existing-branch] was not found in the repository [oc-org/my-app-repo]. Please check configuration."
}