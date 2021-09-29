# Control-4
package gitlab.keys.readonly

default control = "Control 4"
default message = ""

is_read_only[message] {
    message := verify(input.Title, input.ReadOnly)
}

verify(keyTitle, isReadOnly) = result {
	isReadOnly == true
    result := sprintf("%v: INFO - Automation key with name [%v] is correctly set-up as read-only.", [control, keyTitle])
}

verify(keyTitle, isReadOnly) = result {
	isReadOnly == false
    result := sprintf("%v: WARNING - Automation key with name [%v] is not read-only. Please consider updating it to follow principle of least privilege access.",
        [control, keyTitle]
    )
}