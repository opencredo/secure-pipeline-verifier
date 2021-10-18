# Control-3
package token.expiry

default control = "Control 3"

currentTime = time.now_ns()

needs_update = decision {
	keyCreationDate := time.parse_rfc3339_ns(input.CreationDate)

    y := 0
    m := 1
    d := 0
    nextUpdateExpiry := time.add_date(keyCreationDate, y, m, d)

    # if nextUpdateExpiry is lower than currentTime means the key hasn't been updated in more than a month
 	decision := verify(input.Title, currentTime, nextUpdateExpiry)
}

verify(keyTitle, currentDateTime, nextExpectedUpdateDateTime) = decision {
	currentDateTime >= nextExpectedUpdateDateTime
    response := sprintf("Automation key [%v] has not been changed for more than a month. Please consider updating it.", [keyTitle])

    decision := {"control": control, "level": "WARNING", "msg": response}
}

verify(keyTitle, currentDateTime, nextExpectedUpdateDateTime) = decision {
	currentDateTime <= nextExpectedUpdateDateTime
    response := sprintf("Automation key [%v] does not need to be updated at this time.", [keyTitle])

    decision := {"control": control, "level": "INFO", "msg": response}
}