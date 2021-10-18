# Control-4 test
package keys.readonly

test_key_is_readonly {

    safe_key_input := {
                        "ID": "12",
                        "Title": "Deploy-Key-UAT",
                        "ReadOnly": true
                      }

    expected := {
        "control": "Control 4",
        "level": "INFO",
        "msg": "Automation key with name [Deploy-Key-UAT] is correctly set-up as read-only."
    }

    decision := is_read_only with input as safe_key_input

    decision.control == "Control 4"
    decision.level == "INFO"
    decision.msg == "Automation key with name [Deploy-Key-UAT] is correctly set-up as read-only."
}

test_key_is_not_readonly {
    unsafe_key_input := {
                          "ID": "12",
                          "Title": "Deploy-Key-DEV",
                          "ReadOnly": false
                        }

    expected := {
        "control": "Control 4",
        "level": "WARNING",
        "msg": "Automation key with name [Deploy-Key-DEV] is not read-only. Please consider updating it to follow principle of least privilege access."
    }

    decision := is_read_only with input as unsafe_key_input

    decision.control == "Control 4"
    decision.level == "WARNING"
    decision.msg == "Automation key with name [Deploy-Key-DEV] is not read-only. Please consider updating it to follow principle of least privilege access."
}