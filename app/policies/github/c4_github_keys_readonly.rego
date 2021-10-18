# Control-4
package keys.readonly

default control = "Control 4"

is_read_only = decision {
    decision := verify(input.Title, input.ReadOnly)
}

verify(keyTitle, isReadOnly) = decision {
	isReadOnly == true
    response := sprintf("Automation key with name [%v] is correctly set-up as read-only.", [keyTitle])

    decision := {"control": control, "level": "INFO", "msg": response}
}

verify(keyTitle, isReadOnly) = decision {
	isReadOnly == false
    response := sprintf("Automation key with name [%v] is not read-only. Please consider updating it to follow principle of least privilege access.",
        [keyTitle]
    )

    decision := {"control": control, "level": "WARNING", "msg": response}
}