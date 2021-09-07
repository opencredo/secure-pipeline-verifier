# Control-2
package github.branch.protection

default message = ""

is_unprotected[message] {
	input.SignatureProtected == false
	contains(input.Error, "Branch not protected")
    message := sprintf("WARNING - The branch [%v] of repository [%v] is not protected with signed commits as expected. Please consider protecting it.", [input.BranchName, input.GitHubRepo])
}