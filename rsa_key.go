package login_sdk_go

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
)

var (
	ErrNoPublicRSA   = errors.New("there is no public RSA keys available for this project")
	ErrNoKeysMatches = errors.New("unable to find a signing key that matches")
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
		return nil, ErrNoPublicRSA
	}

	for i := range keysResp {
		if keysResp[i].Kid == kid {
			return keysResp[0].CreateRSAPublicKey(), nil
		}
	}

	return nil, fmt.Errorf("%w %s", ErrNoKeysMatches, kid)
}
