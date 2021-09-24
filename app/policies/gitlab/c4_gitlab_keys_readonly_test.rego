# Control-4 test
package gitlab.keys.readonly

test_key_is_readonly {

    safe_key_input := {
                        "ID": "12",
                        "Title": "Deploy-Key-UAT",
                        "ReadOnly": true
                      }

    expected := "INFO - Automation key with name [Deploy-Key-UAT] is correctly set-up as read-only."

    is_read_only[expected] with input as safe_key_input
}

test_key_is_not_readonly {
    unsafe_key_input := {
                          "ID": "12",
                          "Title": "Deploy-Key-DEV",
                          "ReadOnly": false
                        }

    expected := "WARNING - Automation key with name [Deploy-Key-DEV] is not read-only. Please consider updating it to follow principle of least privilege access."

    is_read_only[expected] with input as unsafe_key_input
}