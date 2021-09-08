# Control-4
package github.keys.readonly

default message = ""

is_read_only[message] {
    message := verify(input.Title, input.ReadOnly)
}

verify(keyTitle, isReadOnly) = result {
	isReadOnly == true
    result := sprintf("INFO - Automation key with name [%v] is correctly set-up as read-only.", [keyTitle])
}

verify(keyTitle, isReadOnly) = result {
	isReadOnly == false
    result := sprintf("WARNING - Automation key with name [%v] is not read-only. Please consider updating it to follow principle of least privilege access.", [keyTitle])
}