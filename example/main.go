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
	_, validateErr := loginSdk.Validate(context.Background(), "{YOUR_TOKEN}")
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
