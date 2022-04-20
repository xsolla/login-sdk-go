package login_sdk_go

import (
	"gitlab.loc/sdk-login/login-sdk-go/cache"
	"gitlab.loc/sdk-login/login-sdk-go/model"
)

type ProjectKeysGetter interface {
	GetProjectKeysForLoginProject(projectID string) ([]model.ProjectPublicKey, error)
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

func (c cachedValidationKeysStorage) GetProjectKeysForLoginProject(projectID string) ([]model.ProjectPublicKey, error) {
	cached, found := c.cache.Get(projectID)

	if found {
		return cached.([]model.ProjectPublicKey), nil
	}

	res, err := c.client.GetProjectKeysForLoginProject(projectID)
	if err == nil {
		c.cache.Set(projectID, res)
	}

	return res, err
}
