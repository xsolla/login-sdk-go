package login_sdk_go

import (
	"crypto/rsa"
	"errors"
)

type ValidationKeysGetter interface {
	GetProjectKeysForLoginProject(projectId string) ([]RSAKey, error)
}

type RSASigningKey struct {
	projectId string
	storage   ValidationKeysGetter
}

func (rsa RSASigningKey) getSigningKey(kid string) (*rsa.PublicKey, error) {
	keysResp, err := rsa.storage.GetProjectKeysForLoginProject(rsa.projectId)

	if err != nil {
		return nil, err
	}

	if len(keysResp) == 0 {
		return nil, errors.New("there is no public RSA keys available for this project")
	}

	for i := range keysResp {
		if keysResp[i].Kid == kid {
			return createPublicKey(keysResp[0]), nil
		}
	}

	return nil, errors.New("unable to find a signing key that matches " + kid)
}
