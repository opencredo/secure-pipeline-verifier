# Control-2
package gitlab.branch.protection

default message = ""

is_protected[message] {
    message := verify(input)
}

verify(branchInfo) = response {
    	branchInfo.SignatureProtected == false
        response := sprintf("WARNING - The branch [%v] of repository [%v] is not protected with signed commits as expected. Please consider protecting it.",
            [branchInfo.BranchName, branchInfo.Repo]
        )
}

verify(branchInfo) = response {
    	branchInfo.SignatureProtected == true
        response := sprintf("INFO - The branch [%v] of repository [%v] is protected with signed commits as expected.",
            [branchInfo.BranchName, branchInfo.Repo]
        )
}