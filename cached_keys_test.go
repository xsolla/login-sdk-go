package login_sdk_go

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type cacheTestClient struct {
	result *[]RSAKey
}

type cacheTestCache struct {
	cached *[]RSAKey
}

func (f cacheTestClient) GetProjectKeysForLoginProject(projectId string) ([]RSAKey, error) {
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
		var storedInCache = []RSAKey{
			{
				Alg:      "RS256",
				Exponent: "10001",
				Kid:      "1",
				Kty:      "123",
				Modulus:  "456",
				Use:      "use",
			},
		}

		var givenByClient []RSAKey

		var cks = NewCachedValidationKeysStorage(cacheTestClient{result: &givenByClient}, cacheTestCache{cached: &storedInCache})

		res, err := cks.GetProjectKeysForLoginProject("testId")
		require.NoError(t, err)
		require.Equal(t, storedInCache, res)
	})

	t.Run("get fresh", func(t *testing.T) {
		var storedInCache []RSAKey

		var givenByClient = []RSAKey{
			{
				Alg:      "RS256",
				Exponent: "10001",
				Kid:      "1",
				Kty:      "123",
				Modulus:  "456",
				Use:      "use",
			},
		}

		var cks = NewCachedValidationKeysStorage(cacheTestClient{result: &givenByClient}, cacheTestCache{cached: &storedInCache})

		res, err := cks.GetProjectKeysForLoginProject("new_project")
		require.NoError(t, err)
		require.Equal(t, givenByClient, res)
	})
}
