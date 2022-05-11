package keys

import (
	"context"
	"errors"

	"github.com/dgrijalva/jwt-go"

	"gitlab.loc/sdk-login/login-sdk-go/internal/contract"
)

var ErrSHASecretKeyIsEmpty = errors.New("sha secret key is empty")

type HS256SigningKeyGetter struct {
	key string
}

func NewHS256SigningKeyGetter(key string) HS256SigningKeyGetter {
	return HS256SigningKeyGetter{
		key: key,
	}
}

type RS256SigningKeyGetter struct {
	rsaPublicKeyGetter RSAPublicKeyGetter
}

func NewRS256SigningKeyGetter(keyGetter RSAPublicKeyGetter) RS256SigningKeyGetter {
	return RS256SigningKeyGetter{
		rsaPublicKeyGetter: keyGetter,
	}
}

func (hs HS256SigningKeyGetter) GetKey(context.Context, *jwt.Token) (interface{}, error) {
	if hs.key == "" {
		return nil, ErrSHASecretKeyIsEmpty
	}

	return []byte(hs.key), nil
}

func (rs RS256SigningKeyGetter) GetKey(ctx context.Context, token *jwt.Token) (interface{}, error) {
	if kid, ok := token.Header["kid"].(string); ok {
		cl, ok := token.Claims.(contract.SDKClaims)
		if !ok {
			return nil, errors.New("failed receive claims for token")
		}

		key, err := rs.rsaPublicKeyGetter.getPublicKey(ctx, kid, cl.GetProjectID())
		if err != nil {
			return nil, err
		}

		return key, nil
	}

	return nil, errors.New("token doesn't have kid header")
}
