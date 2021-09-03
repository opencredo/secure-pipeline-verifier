# Control-2
package github.branch.protection

default allow = false

allow {
	input.SignatureProtected == true
	not contains(input.Error, "Branch not protected")
}