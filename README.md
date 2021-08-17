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
