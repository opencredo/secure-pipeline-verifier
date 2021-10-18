# Control-2
package signature.protection

default control = "Control 2"

is_protected = decision {
    decision := verify(input)
}

verify(branchInfo) = decision {
    branchInfo.SignatureProtected == false
    contains(branchInfo.Error, "Branch not protected")
    response := sprintf("The branch [%v] of repository [%v] is not protected with signed commits as expected. Please consider protecting it.",
        [branchInfo.BranchName, branchInfo.Repo]
    )

   decision := {"control": control, "level": "WARNING", "msg": response}
}

verify(branchInfo) = decision {
    branchInfo.SignatureProtected == false
    contains(branchInfo.Error, "Not Found")
    response := sprintf("The user has not Admin permissions on repository [%v] to perform this check. Please consider updating permissions.",
        [branchInfo.Repo]
    )

    decision := {"control": control, "level": "ERROR", "msg": response}
}

verify(branchInfo) = decision {
    branchInfo.SignatureProtected == false
    contains(branchInfo.Error, "Branch not found")
    response := sprintf("The branch [%v] was not found in the repository [%v]. Please check configuration.",
        [branchInfo.BranchName, branchInfo.Repo]
    )

    decision := {"control": control, "level": "ERROR", "msg": response}
}

verify(branchInfo) = decision {
    branchInfo.SignatureProtected == true
    branchInfo.Error == ""
    response := sprintf("The branch [%v] of repository [%v] is protected with signed commits as expected.",
        [branchInfo.BranchName, branchInfo.Repo]
    )

    decision := {"control": control, "level": "INFO", "msg": response}
}