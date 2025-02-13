# Login SDK (Golang)
The Xsolla Login SDK for Go is designed to validate various types of JWTs
used for authorizing requests via Xsolla Login.
The SDK supports token validation with the following encryption algorithms:
* HS253 (HMAC): Utilizes a secret key specified during SDK configuration;
* RS256: Employs a public key to verify the token;

## Usage

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
If validation with HS256 tokens is required, specify the secret as follows:
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
Note: If a secret is specified, it will be exclusively used
for validating HS256 tokens. If the secret is not provided,
the Login HTTP API will be utilized for such tokens, allowing
validation across **different projects**.

#### Validating with Custom Claims:
There are two methods to validate tokens with custom claims:
1. Composing `login_sdk_go.CustomClaims` into your custom claims;
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
2. Fully implementing the `Claims` interface, specifically the `Valid()` and `GetProjectID()` methods (the latter is required for working with the Login API);
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

### Configuration options:
```
ShaSecretKey        - Symmetric key (if symmetric signing is used);
Cache               - Interface for cache operations (optional);
ExtraHeaderName     - Name for an additional header;
ExtraHeaderValue    - Value for the additional header;
APITimeout          - Timeout for API calls. Defaults to 5 seconds (optional).
```
**NOTE:** When using a custom cache implementation (e.g., via Redis),
ensure that the TTL for storing keys does not exceed 10 minutes.
If no custom cache is provided, an in-memory cache will be used.
**Additional Header:**
An additional header is used to bypass rate-limiting mechanisms on the Login API side.
The header's name and value are provided individually by each team.
If the header name is not specified, the header is not added.

#### Example of Redis Cache Implementation:
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

### Example
For usage examples, see the `/example`directory.
To build and run, you can use the following commands:
```shell
> make build 
> make run
```

### Testing & Mocks

For ability testing application (replace Login API response with mockproxy) you can set https proxy address, for it use environment variable `HTTPS_PROXY`

### Possible Issues

A potential loading issue may occur:
```
go get -u "gitlab.loc/sdk-login/login-sdk-go"
go: downloading gitlab.loc/sdk-login/login-sdk-go v0.1.1
go get: gitlab.loc/sdk-login/login-sdk-go@v0.1.1: verifying module: gitlab.loc/sdk-login/login-sdk-go@v0.1.1: reading https://sum.golang.org/lookup/gitlab.loc/sdk-login/login-sdk-go@v0.1.1: 410 Gone server response: not found: gitlab.loc/sdk-login/login-sdk-go@v0.1.1: unrecognized import path "gitlab.loc/sdk-login/login-sdk-go": https fetch: Get "https://gitlab.loc/sdk-login/login-sdk-go?go-get=1": dial tcp: lookup gitlab.loc on 8.8.8.8:53: no such host 
```
Cause: Loading code from a private repository.

Solution: Set the GOPRIVATE variable as follows:
```
go env -w GOPRIVATE=gitlab.loc/sdk-login/login-sdk-go
```
To verify, execute the following:
```
go get -u "gitlab.loc/sdk-login/login-sdk-go"
```