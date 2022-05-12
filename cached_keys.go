package login_sdk_go

import (
	"context"
	"errors"

	"gitlab.loc/sdk-login/login-sdk-go/cache"
	"gitlab.loc/sdk-login/login-sdk-go/model"
)

type ProjectKeysGetter interface {
	GetProjectKeysForLoginProject(ctx context.Context, projectID string) ([]model.ProjectPublicKey, error)
}

var ErrConvertKey = errors.New("error converting to ProjectPublicKey")

type CachedValidationKeysStorage struct {
	client ProjectKeysGetter
	cache  cache.ValidationKeysCache
}

func NewCachedValidationKeysStorage(client ProjectKeysGetter, cache cache.ValidationKeysCache) CachedValidationKeysStorage {
	return CachedValidationKeysStorage{
		client: client,
		cache:  cache,
	}
}

func (c CachedValidationKeysStorage) GetProjectKeysForLoginProject(ctx context.Context, projectID string) ([]model.ProjectPublicKey, error) {
	cached, found := c.cache.Get(projectID)

	if found {
		key, ok := cached.([]model.ProjectPublicKey)
		if !ok {
			return nil, ErrConvertKey
		}

		return key, nil
	}

	res, err := c.client.GetProjectKeysForLoginProject(ctx, projectID)
	if err == nil {
		c.cache.Set(projectID, res)
	}

	return res, err
}
