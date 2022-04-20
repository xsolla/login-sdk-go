package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type ValidationKeysCache interface {
	Get(projectID string) (interface{}, bool)
	Set(projectID string, keys interface{})
}

type DefaultValidationKeysCache struct {
	cache *cache.Cache
}

func NewDefaultCache(expirationTime time.Duration) DefaultValidationKeysCache {
	return DefaultValidationKeysCache{
		cache.New(expirationTime, expirationTime),
	}
}

func (dc DefaultValidationKeysCache) Get(projectID string) (interface{}, bool) {
	return dc.cache.Get(projectID)
}

func (dc DefaultValidationKeysCache) Set(projectID string, keys interface{}) {
	dc.cache.Set(projectID, keys, cache.DefaultExpiration)
}
