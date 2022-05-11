package login_sdk_go

import (
	"context"
	"errors"

	"github.com/dgrijalva/jwt-go"

	sdkcontract "gitlab.loc/sdk-login/login-sdk-go/contract"
	"gitlab.loc/sdk-login/login-sdk-go/internal/contract"
)

var (
	ErrReceiveTokenClaims = errors.New("failed receive claims for token")
	ErrNoKidHeader        = errors.New("token doesn't have kid header")
	ErrTypeAssertion      = errors.New("type assertion error")
)

var errSHASecretKeyIsEmpty = errors.New("sha secret key is empty")

type SigningKeyGetter interface {
	getKey(ctx context.Context, token interface{}) (interface{}, error)
}

type HS256SigningKeyGetter struct {
	key string
}

type RS256SigningKeyGetter struct {
	Config
	rsaPublicKeyGetter RSAPublicKeyGetter
}

func (hs HS256SigningKeyGetter) getKey(context.Context, interface{}) (interface{}, error) {
	if hs.key == "" {
		return nil, errSHASecretKeyIsEmpty
	}

	return []byte(hs.key), nil
}

func (rs RS256SigningKeyGetter) getKey(ctx context.Context, token interface{}) (interface{}, error) {
	jwtToken, ok := token.(*jwt.Token)
	if !ok {
		return nil, ErrTypeAssertion
	}

	if kid, ok := jwtToken.Header["kid"].(string); ok {
		claims, ok := jwtToken.Claims.(sdkcontract.Claims)
		if !ok {
			return nil, ErrReceiveTokenClaims
		}
		rs.rsaPublicKeyGetter.projectID = claims.GetProjectID()
		key, err := rs.rsaPublicKeyGetter.getPublicKey(ctx, kid)
		if err != nil {
			return nil, err
		}

		return key, nil
	}

	return nil, ErrNoKidHeader
}
