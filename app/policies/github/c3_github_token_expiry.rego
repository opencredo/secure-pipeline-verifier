# Control-3
package github.token.expiry

default control = "Control 3"
default message = ""
currentTime = time.now_ns()

needs_update[message] {
	keyCreationDate := time.parse_rfc3339_ns(input.CreationDate)

    y := 0
    m := 1
    d := 0
    nextUpdateExpiry := time.add_date(keyCreationDate, y, m, d)

    # if nextUpdateExpiry is lower than currentTime means the key hasn't been updated in more than a month
 	message := verify(input.Title, currentTime, nextUpdateExpiry)
}

verify(keyTitle, currentDateTime, nextExpectedUpdateDateTime) = result {
	currentDateTime >= nextExpectedUpdateDateTime
    result := sprintf("%v: WARNING - Automation key [%v] has not been changed for more than a month. Please consider updating it.",
        [control, keyTitle]
    )
}

verify(keyTitle, currentDateTime, nextExpectedUpdateDateTime) = result {
	currentDateTime <= nextExpectedUpdateDateTime
    result := sprintf("%v: INFO - Automation key [%v] does not need to be updated at this time.", [control, keyTitle])
}