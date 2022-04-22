package login_sdk_go

import (
	"context"

	"gitlab.loc/sdk-login/login-sdk-go/cache"
	"gitlab.loc/sdk-login/login-sdk-go/model"
)

type ProjectKeysGetter interface {
	GetProjectKeysForLoginProject(ctx context.Context, projectId string) ([]model.ProjectPublicKey, error)
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

func (c cachedValidationKeysStorage) GetProjectKeysForLoginProject(ctx context.Context, projectId string) ([]model.ProjectPublicKey, error) {
	cached, found := c.cache.Get(projectId)

	if found {
		return cached.([]model.ProjectPublicKey), nil
	}

	res, err := c.client.GetProjectKeysForLoginProject(ctx, projectId)
	if err == nil {
		c.cache.Set(projectId, res)
	}

	return res, err
}
