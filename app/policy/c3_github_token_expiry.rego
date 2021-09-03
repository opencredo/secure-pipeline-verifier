# Control-3
package github.token.expiry

default allow = false

allow {
	mockDateNow := time.parse_rfc3339_ns("2021-08-15T00:00:00Z")
	tokenUpdateDate := time.parse_rfc3339_ns(input.updated_at)

    y := 0
    m := 1
    d := 0
    nextUpdateExpiry := time.add_date(tokenUpdateDate, y, m, d)

	mockDateNow < nextUpdateExpiry
}