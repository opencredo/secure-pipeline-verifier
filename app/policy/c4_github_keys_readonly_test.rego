package github.keys.readonly

test_key_is_readonly {
    allow with input as {
                          "ID": "12",
                          "Title": "Deploy-Key-UAT",
                          "Verified": true,
                          "ReadOnly": true
                        }
}

test_key_is_not_readonly {
    not allow with input as {
                          "ID": "12",
                          "Title": "Deploy-Key-DEV",
                          "Verified": true,
                          "ReadOnly": false
                        }
}