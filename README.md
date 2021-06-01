## Usage

```go
package main

import (
    "gitlab.loc/sdk-login/login-sdk-go"
)


func main() {
    var config = login_sdk_go.Config{
        ShaSecretKey:      "SECRET",
        LoginProjectId:    "PROJECT_ID",
        LoginClientId:     1,
        LoginClientSecret: "LOGIN_CLIENT_SECRET",
    }
    
    var loginSgk = login_sdk_go.New(&config)

    /* Validate token */
	token, err := loginSgk.Validate("<MY_TOKEN>")
	
    /* Refreshing token */
	response, err := loginSgk.Refresh("<REFRESH_TOKEN>")
}