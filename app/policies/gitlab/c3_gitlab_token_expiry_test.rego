# Control-3 test
package token.expiry

test_key_needs_not_update {
    # Sept 6th 2021 - 15:32
    mock_currentTime := 1630938710180067000

    safe_key_input := {
                  "ID": 56618754,
                  "Title": "test-deploy-key",
                  "ReadOnly": true,
                  "CreationDate": "2021-09-06T13:31:13Z"
                }

    expected := {
        "control": "Control 3",
        "level": "INFO",
        "msg": "Automation key [test-deploy-key] does not need to be updated at this time."
    }

    decision := needs_update with currentTime as mock_currentTime with input as safe_key_input

    decision.control == "Control 3"
    decision.level == "INFO"
    decision.msg == "Automation key [test-deploy-key] does not need to be updated at this time."
}

test_key_needs_update {
    # Sept 6th 2021 - 15:32
    mock_currentTime := 1630938710180067000

    unsafe_key_input := {
                      "ID": 56618754,
                      "Title": "my-old-deploy-key",
                      "ReadOnly": true,
                      "CreationDate": "2021-08-03T15:21:56Z"
                  }

    expected := {
        "control": "Control 3",
        "level": "WARNING",
        "msg": "Automation key [my-old-deploy-key] has not been changed for more than a month. Please consider updating it."
    }

    decision := needs_update with currentTime as mock_currentTime with input as unsafe_key_input

    decision.control == "Control 3"
    decision.level == "WARNING"
    decision.msg == "Automation key [my-old-deploy-key] has not been changed for more than a month. Please consider updating it."
}