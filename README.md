## Usage

```go
package main

import (
    "gitlab.loc/sdk-login/login-sdk-go"
)


func main() {
    config := login_sdk_go.Config{
        ShaSecretKey:      "SECRET",
        LoginProjectId:    "PROJECT_ID",
        LoginClientId:     1,
        LoginClientSecret: "LOGIN_CLIENT_SECRET",
    }
    
    loginSdk := login_sdk_go.New(&config)
    
	err := loginSgk.Validate("<MY_TOKEN>")
	
	if login_sdk_go.IsExpiredErr(err) {
		newToken := loginSdk.Refresh("Sdsds")
	}
}