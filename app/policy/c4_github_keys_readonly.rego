# Control-4
package github.keys.readonly

default message = ""

can_write[message] {
	input.ReadOnly == false
    message := sprintf("WARNING - Automation key with name [%v] is not read-only. Please consider updating it to follow principle of least privilege access.", [input.Title])
}