package login_sdk_go

import (
	"context"
	"crypto/rsa"
	"errors"
)

type RSAPublicKeyGetter struct {
	projectID string
	storage   ProjectKeysGetter
}

func (rsa RSAPublicKeyGetter) getPublicKey(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	keysResp, err := rsa.storage.GetProjectKeysForLoginProject(ctx, rsa.projectID)

	if err != nil {
		return nil, err
	}

	if len(keysResp) == 0 {
		return nil, errors.New("there is no public RSA keys available for this project")
	}

	for i := range keysResp {
		if keysResp[i].Kid == kid {
			return keysResp[0].CreateRSAPublicKey(), nil
		}
	}

	return nil, errors.New("unable to find a signing key that matches " + kid)
}
