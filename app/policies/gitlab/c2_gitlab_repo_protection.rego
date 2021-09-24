# Control-2
package gitlab.repo.protection

default message = ""

is_protected[message] {
    message := verify(input)
}

verify(repoInfo) = response {
    	repoInfo.SignatureProtected == false
        response := sprintf("WARNING - The repository [%v] is not protected with signed commits as expected. Please consider protecting it.",
            [repoInfo.Repo]
        )
}

verify(repoInfo) = response {
    	repoInfo.SignatureProtected == true
        response := sprintf("INFO - The repository [%v] is protected with signed commits as expected.", [repoInfo.Repo])
}