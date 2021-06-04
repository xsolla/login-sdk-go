package login_sdk_go

import (
	"gitlab.loc/sdk-login/login-sdk-go/cache"
	"gitlab.loc/sdk-login/login-sdk-go/model"
)

type ProjectKeysGetter interface {
	GetProjectKeysForLoginProject(projectId string) ([]model.ProjectPublicKey, error)
}

type cachedValidationKeysStorage struct {
	client ProjectKeysGetter
	cache  cache.ValidationKeysCache
}

func NewCachedValidationKeysStorage(client ProjectKeysGetter, cache cache.ValidationKeysCache) ProjectKeysGetter {
	return cachedValidationKeysStorage{
		client: client,
		cache:  cache,
	}
}

func (c cachedValidationKeysStorage) GetProjectKeysForLoginProject(projectId string) ([]model.ProjectPublicKey, error) {
	cached, found := c.cache.Get(projectId)

	if found {
		return cached.([]model.ProjectPublicKey), nil
	}

	res, err := c.client.GetProjectKeysForLoginProject(projectId)
	c.cache.Set(projectId, res)

	cached, found = c.cache.Get(projectId)

	return res, err
}
