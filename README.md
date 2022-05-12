# Login SDK (Golang)
Данный SDK предназначет для валидации различных видов JWT, которые 
используются для авторизации запросов с помощью Xsolla Login.
SDK поддерживает валидацию токенов со следующими алгоритмами шифрования:
* HS253 (HMAC): используется секрет, задаваемый при конфигурации SDK;
* RS256: используется публичный ключ для проверки токена;

## Использование 

```go
package main

import (
    "gitlab.loc/sdk-login/login-sdk-go"
    "fmt"
    "os"
)


func main() {
    config := login_sdk_go.Config{}
    
    loginSdk, err := login_sdk_go.New(config)
    
    if err != nil {
        fmt.Println("Failed inititalize SDK: %s", err)
        os.Exit(1)
    }
	
    _, err = loginSdk.Validate("<MY_TOKEN>")

    if !err.Valid() {
    	fmt.Println("Error: ", err.Error())
    	os.Exit(1)
    }
    fmt.Println("Success!")
}
```
При необходимости такой валидации необходимо указать секрет:
```go
package main

import (
    "gitlab.loc/sdk-login/login-sdk-go"
    "fmt"
    "os"
)

func main() {
    config := login_sdk_go.Config{
		ShaSecretKey: "your-secret",
    }
    
    loginSdk, err := login_sdk_go.New(config)
    
    if err != nil {
        fmt.Println("Failed inititalize SDK: %s", err)
        os.Exit(1)
    }
    _, err = loginSdk.Validate("<MY_TOKEN>")

    if !err.Valid() {
        fmt.Println("Error: ", err.Error())
        os.Exit(1)
    }
    fmt.Println("Success!")
}
```
#### Для валидации с кастомными(своими) клеймсами есть 2 способа:
1. Композиция `login_sdk_go.CustomClaims` в свои кастомные клеймсы;
```go
package main

import (
	"fmt"
	"log"

	"gitlab.loc/sdk-login/login-sdk-go"
)

type YourClaims struct {
	MyField string `json:"my_field"`
	login_sdk_go.CustomClaims
}

func main() {
	config := login_sdk_go.Config{}

	loginSdk, err := login_sdk_go.New(config)

	if err != nil {
		log.Fatalf("Failed inititalize SDK: %v", err)
	}

	var claims YourClaims
	_, vErr := loginSdk.ValidateWithClaims("<MY_TOKEN>", &claims)
	// _, vErr := loginSdk.ValidateWithClaimsAndContext(context.Background(), "<MY_TOKEN>", &claims)

	if !vErr.Valid() {
		log.Fatalf("Validation error: %v", vErr)
	}

	fmt.Println("Success!")
}
```
2. Полная реализация интерфейса `Claims`, а именно, методы `Valid()` и `GetProjectID()` последний нужен для работы с LoginAPI.
```go
package main

import (
	"fmt"
	"log"

	"gitlab.loc/sdk-login/login-sdk-go"
)

type YourClaims struct {
	ProjectID string `json:"xsolla_project_id"`
}

func (c YourClaims) Valid() error {
	return nil
}

func (c YourClaims) GetProjectID() string {
	return c.ProjectID
}

func main() {
	config := login_sdk_go.Config{}

	loginSdk, err := login_sdk_go.New(config)

	if err != nil {
		log.Fatalf("Failed inititalize SDK: %v", err)
	}

	var claims YourClaims
	_, vErr := loginSdk.ValidateWithClaims("<MY_TOKEN>", &claims)
	// _, vErr := loginSdk.ValidateWithClaimsAndContext(context.Background(), "<MY_TOKEN>", &claims)

	if !vErr.Valid() {
		log.Fatalf("Validation error: %v", vErr)
	}

	fmt.Println("Success!")
}
```

Следует отметить, что если указан секрет, то только он будет использоваться
для валидации токенов HS256. В случае, если секрет не указан, то для таких токенов
будет использоваться Login HTTP API, т.е. таким образом можно валидировать токены
**разных проектов**.

#### Для валидации с кастомными(своими) клеймсами есть 2 способа:
1. Композиция `login_sdk_go.CustomClaims` в свои кастомные клеймсы;
```go
package main

import (
	"fmt"
	"log"

	"gitlab.loc/sdk-login/login-sdk-go"
)

type YourClaims struct {
	MyField string `json:"my_field"`
	login_sdk_go.CustomClaims
}

func main() {
	config := login_sdk_go.Config{}

	loginSdk, err := login_sdk_go.New(config)

	if err != nil {
		log.Fatalf("Failed inititalize SDK: %v", err)
	}

	var claims YourClaims
	_, vErr := loginSdk.ValidateWithClaims("<MY_TOKEN>", &claims)
	// _, vErr := loginSdk.ValidateWithClaimsAndContext(context.Background(), "<MY_TOKEN>", &claims)

	if !vErr.Valid() {
		log.Fatalf("Validation error: %v", vErr)
	}

	fmt.Println("Success!")
}
```
2. Полная реализация интерфейса `Claims`, а именно, методы `Valid()` и `GetProjectID()` последний нужен для работы с LoginAPI.
```go
package main

import (
	"fmt"
	"log"

	"gitlab.loc/sdk-login/login-sdk-go"
)

type YourClaims struct {
	ProjectID string `json:"xsolla_project_id"`
}

func (c YourClaims) Valid() error {
	return nil
}

func (c YourClaims) GetProjectID() string {
	return c.ProjectID
}

func main() {
	config := login_sdk_go.Config{}

	loginSdk, err := login_sdk_go.New(config)

	if err != nil {
		log.Fatalf("Failed inititalize SDK: %v", err)
	}

	var claims YourClaims
	_, vErr := loginSdk.ValidateWithClaims("<MY_TOKEN>", &claims)
	// _, vErr := loginSdk.ValidateWithClaimsAndContext(context.Background(), "<MY_TOKEN>", &claims)

	if !vErr.Valid() {
		log.Fatalf("Validation error: %v", vErr)
	}

	fmt.Println("Success!")
}
```

### Конфигурация:
```
ShaSecretKey        - Симметричный ключ (если используется симметричная подпись)
Cache               - Интерфейс для работы с кэшом (опционально).
```
**NOTE:** При использовании кастомной реализации кеша (например через Redis), необходимо
обязательно указывать TTL для хранения ключей не более чем в 10 минут. <br>
В случае, если нет кастомной реализации кеша, будет использоваться кеш, располагающийся в 
оперативной памяти приложения.

#### Пример реализации кеша через Redis:
```go
package cache

import (
	"github.com/go-redis/redis"
	"time"
)

const (
	keyTTL = 10 * time.Minute
)

type RedisCache struct {
	client *redis.Client
	s      *login_sdk_go.ProjectKeysGetter
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

func (r RedisCache) Get(projectId string) (interface{}, bool) {
	if r.client == nil {
		return "", false
	}
	result, err := r.client.Get(projectId).Bytes()
	var res []model.ProjectPublicKey
	err = json.Unmarshal(result, &res)

	return res, err == nil
}

func (r RedisCache) Set(projectId string, value interface{}) {
	b, _ := json.Marshal(value)
	r.client.Set(projectId, b, 10*time.Minute)
}
```

### Пример
Пример использования см в `/example`<br>
Для сборки и запуска можно использовать следующие команды:
```shell
> make build 
> make run
```

### Возможные проблемы

Возможна проблема загрузки:
```
go get -u "gitlab.loc/sdk-login/login-sdk-go"
go: downloading gitlab.loc/sdk-login/login-sdk-go v0.1.1
go get: gitlab.loc/sdk-login/login-sdk-go@v0.1.1: verifying module: gitlab.loc/sdk-login/login-sdk-go@v0.1.1: reading https://sum.golang.org/lookup/gitlab.loc/sdk-login/login-sdk-go@v0.1.1: 410 Gone server response: not found: gitlab.loc/sdk-login/login-sdk-go@v0.1.1: unrecognized import path "gitlab.loc/sdk-login/login-sdk-go": https fetch: Get "https://gitlab.loc/sdk-login/login-sdk-go?go-get=1": dial tcp: lookup gitlab.loc on 8.8.8.8:53: no such host 
```
Причина: загрузка кода из приватного репозитория. <br>
Для решения необходимо установить переменную GOPRIVATE следующим образом:
```
go env -w GOPRIVATE=gitlab.loc/sdk-login/login-sdk-go
```
Для проверки необходимо выполнить: 
```
go get -u "gitlab.loc/sdk-login/login-sdk-go"
```