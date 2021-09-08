# Control-3 test
package github.token.expiry

test_key_needs_not_update {
    # Sept 6th 2021 - 15:32
    mock_currentTime := 1630938710180067000

    safe_key_input := {
                  "ID": 56618754,
                  "Title": "test-deploy-key",
                  "Verified": true,
                  "ReadOnly": true,
                  "CreationDate": "2021-09-06T13:31:13Z"
                }

    expected := "INFO - Automation key [test-deploy-key] does not need to be updated at this time."

    needs_update[expected] with currentTime as mock_currentTime with input as safe_key_input
}

test_key_needs_update {
    # Sept 6th 2021 - 15:32
    mock_currentTime := 1630938710180067000

    unsafe_key_input := {
                      "ID": 56618754,
                      "Title": "my-old-deploy-key",
                      "Verified": true,
                      "ReadOnly": true,
                      "CreationDate": "2021-08-03T15:21:56Z"
                  }

    expected := "WARNING - Automation key [my-old-deploy-key] has not been changed for more than a month. Please consider updating it."

    needs_update[expected] with currentTime as mock_currentTime with input as unsafe_key_input
}