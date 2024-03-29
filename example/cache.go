package main

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"

	"github.com/xsolla/login-sdk-go/keys"
	"github.com/xsolla/login-sdk-go/model"
)

type RedisCache struct {
	client *redis.Client
	s      *keys.ProjectKeysGetter
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

	if err != nil {
		return nil, false
	}
	var res []model.ProjectPublicKey
	err = json.Unmarshal(result, &res)

	return res, err == nil
}

func (r RedisCache) Set(projectId string, value interface{}) {
	b, err := json.Marshal(value)
	if err != nil {
		return
	}
	r.client.Set(projectId, b, 10*time.Minute)
}
