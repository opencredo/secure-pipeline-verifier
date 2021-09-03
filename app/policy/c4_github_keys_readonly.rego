# Control-4
package github.keys.readonly

default allow = false

allow {
	input.Verified == true
	input.ReadOnly == true
}