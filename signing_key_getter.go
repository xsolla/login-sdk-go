package login_sdk_go

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

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
		return nil, errors.New("ShaSecretKey config is required for sync validation")
	}

	return []byte(hs.key), nil
}

func (rs RS256SigningKeyGetter) getKey(token interface{}) (interface{}, error) {
	jwtToken, ok := token.(*jwt.Token)
	if !ok {
		return nil, errors.New("type assertion error")
	}

	if kid, ok := jwtToken.Header["kid"].(string); ok {
		key, err := rs.rsaPublicKeyGetter.getPublicKey(kid)
		if err != nil {
			return nil, err
		}
		return key, nil
	}
	return nil, errors.New("token doesn't have kid header")
}
