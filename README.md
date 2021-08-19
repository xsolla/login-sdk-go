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
        LoginClientId:     42,
        LoginClientSecret: "LOGIN_CLIENT_SECRET",
        SessionApiHost:    "GRPC_HOST",
        SessionApiPort:    "GRPC_PORT"
    }
    
    loginSdk, err := login_sdk_go.New(&config)
    
    if err != nil {
        os.Exit(1)
    }
    
	err = loginSgk.Validate("<MY_TOKEN>")
	
	if login_sdk_go.IsExpiredErr(err) {
		newToken := loginSdk.Refresh("REFRESH_TOKEN")
	}
}
```
### Конфигурация:

```
LoginClientId       - ID OAuth2 клиента
LoginClientSecret   - Секрет OAuth2 клиента
LoginProjectId      - ID Login проекта
SessionApiHost      - Хост Login Session API (gRPC)
SessionApiPort      - Порт Login Session API (gRPC)
ShaSecretKey        - Симметричный ключ (если используется симметричная подпись)
Cache               - Интерфейс для работы с кэшом (опционально)
```

Пример использования см в "/example/main.go"

### Possible Issues

You can face with error like below:
```
go get -u "gitlab.loc/sdk-login/login-sdk-go"
go: downloading gitlab.loc/sdk-login/login-sdk-go v0.1.1
go get: gitlab.loc/sdk-login/login-sdk-go@v0.1.1: verifying module: gitlab.loc/sdk-login/login-sdk-go@v0.1.1: reading https://sum.golang.org/lookup/gitlab.loc/sdk-login/login-sdk-go@v0.1.1: 410 Gone server response: not found: gitlab.loc/sdk-login/login-sdk-go@v0.1.1: unrecognized import path "gitlab.loc/sdk-login/login-sdk-go": https fetch: Get "https://gitlab.loc/sdk-login/login-sdk-go?go-get=1": dial tcp: lookup gitlab.loc on 8.8.8.8:53: no such host 
```

RootCause: U try download from private repository

For resolve that issue just set GOPRIVATE as show below:
```
go env -w GOPRIVATE=gitlab.loc/sdk-login/login-sdk-go
```
And test that
```
go get -u "gitlab.loc/sdk-login/login-sdk-go"
```