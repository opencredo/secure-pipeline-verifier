# Control-2
package github.branch.protection

default message = ""

is_protected[message] {
    message := verify(input)
}

verify(branchInfo) = response {
    	branchInfo.SignatureProtected == false
    	contains(branchInfo.Error, "Branch not protected")
        response := sprintf("WARNING - The branch [%v] of repository [%v] is not protected with signed commits as expected. Please consider protecting it.", [branchInfo.BranchName, branchInfo.GitHubRepo])
}

verify(branchInfo) = response {
    	branchInfo.SignatureProtected == false
    	contains(branchInfo.Error, "Not Found")
        response := sprintf("ERROR - The user has not Admin permissions on repository [%v] to perform this check. Please consider updating permissions.", [branchInfo.GitHubRepo])
}

verify(branchInfo) = response {
    	branchInfo.SignatureProtected == false
    	contains(branchInfo.Error, "Branch not found")
        response := sprintf("ERROR - The branch [%v] was not found in the repository [%v].", [branchInfo.BranchName, branchInfo.GitHubRepo])
}

verify(branchInfo) = response {
    	branchInfo.SignatureProtected == true
    	branchInfo.Error == ""
        response := sprintf("INFO - The branch [%v] of repository [%v] is protected with signed commits as expected.", [branchInfo.BranchName, branchInfo.GitHubRepo])
}