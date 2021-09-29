# Control-2
package github.branch.protection

default control = "Control 2"
default message = ""

is_protected[message] {
    message := verify(input)
}

verify(branchInfo) = response {
    	branchInfo.SignatureProtected == false
    	contains(branchInfo.Error, "Branch not protected")
        response := sprintf("%v: WARNING - The branch [%v] of repository [%v] is not protected with signed commits as expected. Please consider protecting it.",
            [control, branchInfo.BranchName, branchInfo.GitHubRepo]
        )
}

verify(branchInfo) = response {
    	branchInfo.SignatureProtected == false
    	contains(branchInfo.Error, "Not Found")
        response := sprintf("%v: ERROR - The user has not Admin permissions on repository [%v] to perform this check. Please consider updating permissions.",
            [control, branchInfo.GitHubRepo]
        )
}

verify(branchInfo) = response {
    	branchInfo.SignatureProtected == false
    	contains(branchInfo.Error, "Branch not found")
        response := sprintf("%v: ERROR - The branch [%v] was not found in the repository [%v]. Please check configuration.",
            [control, branchInfo.BranchName, branchInfo.GitHubRepo]
        )
}

verify(branchInfo) = response {
    	branchInfo.SignatureProtected == true
    	branchInfo.Error == ""
        response := sprintf("%v: INFO - The branch [%v] of repository [%v] is protected with signed commits as expected.",
            [control, branchInfo.BranchName, branchInfo.GitHubRepo]
        )
}