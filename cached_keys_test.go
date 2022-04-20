package login_sdk_go

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.loc/sdk-login/login-sdk-go/model"
)

type cacheTestClient struct {
	result *[]model.ProjectPublicKey
}

type cacheTestCache struct {
	cached *[]model.ProjectPublicKey
}

func (f cacheTestClient) GetProjectKeysForLoginProject(projectId string) ([]model.ProjectPublicKey, error) {
	return *f.result, nil
}

func (c cacheTestCache) Get(key string) (interface{}, bool) {
	if key == "new_project" {
		return nil, false
	}
	return *c.cached, true
}

func (c cacheTestCache) Set(key string, value interface{}) {}

func TestCachedValidationKeysStorage(t *testing.T) {
	t.Run("get from cache", func(t *testing.T) {
		storedInCache := []model.ProjectPublicKey{
			{
				Alg:      "RS256",
				Exponent: "10001",
				Kid:      "1",
				Kty:      "123",
				Modulus:  "456",
				Use:      "use",
			},
		}

		var givenByClient []model.ProjectPublicKey

		cks := NewCachedValidationKeysStorage(cacheTestClient{result: &givenByClient}, cacheTestCache{cached: &storedInCache})

		res, err := cks.GetProjectKeysForLoginProject("testId")
		require.NoError(t, err)
		require.Equal(t, storedInCache, res)
	})

	t.Run("get fresh", func(t *testing.T) {
		var storedInCache []model.ProjectPublicKey

		givenByClient := []model.ProjectPublicKey{
			{
				Alg:      "RS256",
				Exponent: "10001",
				Kid:      "1",
				Kty:      "123",
				Modulus:  "456",
				Use:      "use",
			},
		}

		cks := NewCachedValidationKeysStorage(cacheTestClient{result: &givenByClient}, cacheTestCache{cached: &storedInCache})

		res, err := cks.GetProjectKeysForLoginProject("new_project")
		require.NoError(t, err)
		require.Equal(t, givenByClient, res)
	})
}
