# Control-4 test
package github.keys.readonly

test_key_is_readonly {
    can_write with input as {
                          "ID": "12",
                          "Title": "Deploy-Key-UAT",
                          "Verified": true,
                          "ReadOnly": true
                        }
}

test_key_is_not_readonly {
    unsafe_key_input := {
                          "ID": "12",
                          "Title": "Deploy-Key-DEV",
                          "Verified": true,
                          "ReadOnly": false
                        }

    expected := "WARNING - Automation key with name [Deploy-Key-DEV] is not read-only. Please consider updating it to follow principle of least privilege access."

    can_write[expected] with input as unsafe_key_input
}