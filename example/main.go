package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"gitlab.loc/sdk-login/login-sdk-go"
	"gitlab.loc/sdk-login/login-sdk-go/model"
	"log"
	"time"
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
	r.client.Set(projectId, b, time.Hour)
}

func main() {
	redisDSN := fmt.Sprintf("%s:%d", "127.0.0.1", 6379)

	redisClient := redis.NewClient(&redis.Options{
		Addr:        redisDSN,
		ReadTimeout: 1 * time.Second,
		PoolSize:    20,
	})

	redisCache := NewRedisCache(redisClient)
	config := login_sdk_go.Config{
		ShaSecretKey:   "test",
		LoginProjectId: "login-id",
		Cache:          redisCache,
	}

	loginSdk, _ := login_sdk_go.New(config)

	parsedToken, err := loginSdk.Validate("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.5mhBHqs5_DTLdINd9p5m7ZJ6XD0Xc55kIaCRY5r6HRA")
	if !err.Valid() {
		if err.IsExpired() {
			// refresh
		}
		log.Fatal(err)
	}

	fmt.Print(parsedToken)
}
