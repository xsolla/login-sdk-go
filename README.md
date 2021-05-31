## Get started

```
package main

import (
    "gitlab.loc/sdk-login/login-sdk-go"
)


func main() {
    var config = login_sdk_go.Config{
        ShaSecretKey:      "SECRET",
        LoginProjectId:    "PROJECT_ID",
        LoginClientId:     "CLIENT_ID",
        LoginClientSecret: "LOGIN_CLIENT_SECRET",
    }
    
    var loginSgk = login_sdk_go.New(&config)
    
    ...
}


```

### Validation

```
...

token, err := loginSgk.Validate("<MY_TOKEN>")
```

### Refreshing

```
...

response, err := loginSgk.Refresh("<REFRESH_TOKEN>")
```