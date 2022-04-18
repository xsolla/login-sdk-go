package login_sdk_go

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

var errSHASecretKeyIsEmpty = errors.New("sha secret key is empty")

type SigningKeyGetter interface {
	getKey(token interface{}) (interface{}, error)
}

type HS256SigningKeyGetter struct {
	key string
}

type RS256SigningKeyGetter struct {
	Config
	rsaPublicKeyGetter RSAPublicKeyGetter
}

func (hs HS256SigningKeyGetter) getKey(interface{}) (interface{}, error) {
	if hs.key == "" {
		return nil, errSHASecretKeyIsEmpty
	}

	return []byte(hs.key), nil
}

func (rs RS256SigningKeyGetter) getKey(token interface{}) (interface{}, error) {
	jwtToken, ok := token.(*jwt.Token)
	if !ok {
		return nil, errors.New("type assertion error")
	}

	if kid, ok := jwtToken.Header["kid"].(string); ok {
		claims, ok := jwtToken.Claims.(*CustomClaims)
		if !ok {
			return nil, errors.New("failed receive claims for token")
		}
		rs.rsaPublicKeyGetter.projectId = claims.ProjectId
		key, err := rs.rsaPublicKeyGetter.getPublicKey(kid)
		if err != nil {
			return nil, err
		}
		return key, nil
	}
	return nil, errors.New("token doesn't have kid header")
}
