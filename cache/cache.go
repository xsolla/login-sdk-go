package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type ValidationKeysCache interface {
	Get(projectId string) (interface{}, bool)
	Set(projectId string, keys interface{})
}

type DefaultValidationKeysCache struct {
	cache *cache.Cache
}

func NewDefaultCache(expirationTime time.Duration) DefaultValidationKeysCache {
	return DefaultValidationKeysCache{
		cache.New(expirationTime, expirationTime),
	}
}

func (dc DefaultValidationKeysCache) Get(projectId string) (interface{}, bool) {
	return dc.cache.Get(projectId)
}

func (dc DefaultValidationKeysCache) Set(projectId string, keys interface{}) {
	dc.cache.Set(projectId, keys, cache.DefaultExpiration)
}
