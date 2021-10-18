# Control-2
package signature.protection

default control = "Control 2"

is_protected = decision {
    decision := verify(input)
}

verify(repoInfo) = decision {
    repoInfo.SignatureProtected == false
    response := sprintf("The repository [%v] is not protected with signed commits as expected. Please consider protecting it.",
        [repoInfo.Repo]
    )

    decision := {"control": control, "level": "WARNING", "msg": response}
}

verify(repoInfo) = decision {
    repoInfo.SignatureProtected == true
    response := sprintf("The repository [%v] is protected with signed commits as expected.",[repoInfo.Repo])

    decision := {"control": control, "level": "INFO", "msg": response}
}