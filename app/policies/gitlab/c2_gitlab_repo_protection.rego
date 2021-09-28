# Control-2
package gitlab.repo.protection

default control = "Control 2"
default message = ""

is_protected[message] {
    message := verify(input)
}

verify(repoInfo) = response {
    	repoInfo.SignatureProtected == false
        response := sprintf("%v: WARNING - The repository [%v] is not protected with signed commits as expected. Please consider protecting it.",
            [control, repoInfo.Repo]
        )
}

verify(repoInfo) = response {
    	repoInfo.SignatureProtected == true
        response := sprintf("%v: INFO - The repository [%v] is protected with signed commits as expected.",
            [control, repoInfo.Repo]
        )
}