package login_sdk_go

func NewCachedValidationKeysStorage(client ValidationKeysGetter, cache Cache) ValidationKeysGetter {
	return cachedValidationKeysStorage{
		client: client,
		cache:  cache,
	}
}

type cachedValidationKeysStorage struct {
	client ValidationKeysGetter
	cache  Cache
}

func (c cachedValidationKeysStorage) GetProjectKeysForLoginProject(projectId string) ([]RSAKey, error) {
	cached, found := c.cache.Get(projectId)

	if found {
		return cached.([]RSAKey), nil
	}

	res, err := c.client.GetProjectKeysForLoginProject(projectId)
	c.cache.Set(projectId, res)

	cached, found = c.cache.Get(projectId)

	return res, err
}
