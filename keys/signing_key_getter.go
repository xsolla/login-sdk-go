package keys

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"

	"gitlab.loc/sdk-login/login-sdk-go/contract"
)

var (
	ErrReceiveTokenClaims  = errors.New("failed receive claims for token")
	ErrNoKidHeader         = errors.New("token doesn't have kid header")
	ErrSHASecretKeyIsEmpty = errors.New("sha secret key is empty")
)

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
		claims, ok := token.Claims.(contract.Claims)
		if !ok {
			return nil, ErrReceiveTokenClaims
		}

		key, err := rs.rsaPublicKeyGetter.getPublicKey(ctx, kid, claims.GetProjectID())
		if err != nil {
			return nil, err
		}

		return key, nil
	}

	return nil, ErrNoKidHeader
}
