package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"

	"gitlab.loc/sdk-login/login-sdk-go"
)

const (
	useRedisCache = true
	shaSecretKey  = ""
)

func main() {

	config := login_sdk_go.Config{}

	if useRedisCache {
		config.Cache = createRedisCache()
	}
	if shaSecretKey != "" {
		config.ShaSecretKey = shaSecretKey
	}

	loginSdk, err := login_sdk_go.New(config)
	if err != nil {
		fmt.Printf("Failed init login sdk. Error: %s", err)
		os.Exit(1)
	}

	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6InNnRnk0NjRrVk5YVFo2YmVYM0tFT2kyam1yWnA4bUQiLCJ0eXAiOiJKV1QifQ.eyJlbWFpbCI6ImdvbmFtYTU2MzJAc3VwZXJ5cC5jb20iLCJleHAiOjE2NTEzMTkyNDksImdyb3VwcyI6W10sImlhdCI6MTY1MDQ1NTI0OSwiaXNfbWFzdGVyIjp0cnVlLCJpc3MiOiJodHRwczovL2xvZ2luLnhzb2xsYS5jb20iLCJwYXJ0bmVyX2RhdGEiOnsiYWRtaW4iOmZhbHNlLCJtZXJjaGFudHMiOlt7ImlkIjoyNzI4NTQsInJvbGUiOiJST0xFX09XTkVSIn1dfSwicGF5bG9hZCI6IntcImZpbHRlcnNcIjpudWxsfSIsInByb21vX2VtYWlsX2FncmVlbWVudCI6ZmFsc2UsInJlZGlyZWN0X3VybCI6Imh0dHBzOi8vcHVibGlzaGVyLnhzb2xsYS5jb20iLCJzdWIiOiI5OTg1MDE5MC05MzBmLTQ4ZGMtYjNhYS0zMThhNTBhODFhMDAiLCJ0eXBlIjoieHNvbGxhX2xvZ2luIiwidXNlcm5hbWUiOiJnb25hbWE1NjMyQHN1cGVyeXAuY29tIiwieHNvbGxhX2xvZ2luX2FjY2Vzc19rZXkiOiJmdGZ4UFJjOTZMejZZVXNJdXdCMmNWV0tXNGpUeTRES3NPQ0NVWnl0WmJZIiwieHNvbGxhX2xvZ2luX3Byb2plY3RfaWQiOiI0MGRiMmVhNC01ZDQyLTExZTYtYTNmZi0wMDUwNTZhMGUwNGEifQ.ky26IpKlCHpcMxFWtMAwrOaSkUGHfvZhsD4mgxAYTP1KXBDdoaJ8icCeywIwmRcIxQZNtrI1R2jekkicSSxOFd20cksW_4ViMwzfKBJgtKZ7zYjUrhXBV8_M5UUjIxclfJrbRukO0ZWSAKvduqvq0c-I5ONylCGxDt0a0o1KQYANAKRogegZb_k3Mn2F3O8Q8a57KHbDSvwQz7oBlLxzQ2-U21alENYERVAuG0dPIQUqpho1XV7XRUUWTZi5EdMdU07EHREg_va98hxJ-C47v4-8qCXTW0iH7anSGmm1-ZlJAQSoGkctMYsQmxCuUpoS5fo0NZ62PAp4hscfwJ_mBw"
	_, validateErr := loginSdk.Validate(context.Background(), token)
	//_, validateErr := loginSdk.Validate(context.Background(), "{YOUR_TOKEN}")
	if !validateErr.Valid() {
		fmt.Println("Error: ", validateErr.Error())
		os.Exit(1)
	}
	fmt.Println("Success!")
}

func createRedisCache() *RedisCache {
	redisDSN := fmt.Sprintf("%s:%d", "127.0.0.1", 6379)

	redisClient := redis.NewClient(&redis.Options{
		Addr:        redisDSN,
		ReadTimeout: 1 * time.Second,
		PoolSize:    20,
	})
	return NewRedisCache(redisClient)
}
