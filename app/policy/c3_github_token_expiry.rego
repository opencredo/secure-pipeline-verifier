# Control-3
package github.token.expiry

default message = ""
currentTime = time.now_ns()

needs_update[message] {
	keyCreationDate := time.parse_rfc3339_ns(input.CreationDate)

    y := 0
    m := 1
    d := 0
    nextUpdateExpiry := time.add_date(keyCreationDate, y, m, d)

    # if nextUpdateExpiry is lower than currentTime means the key hasn't been updated in more than a month
	currentTime >= nextUpdateExpiry
    message := sprintf("WARNING - Automation key [%v] has not been changed since [%v]. Please consider updating it.", [input.Title, input.CreationDate])
}