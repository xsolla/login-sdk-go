package login_sdk_go

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

type SigningKeyGetter interface {
	getKey(token interface{}) (interface{}, error)
}

type HS256SigningKey struct {
	key string
}

type RS256SigningKey struct {
	Config
	publicKeyGetter RSAPublicKey
}

func (hs HS256SigningKey) getKey(interface{}) (interface{}, error) {
	if hs.key == "" {
		return nil, errors.New("ShaSecretKey config is required for sync validation")
	}

	return []byte(hs.key), nil
}

func (rs RS256SigningKey) getKey(token interface{}) (interface{}, error) {
	jwtToken, ok := token.(*jwt.Token)
	if ok == false {
		return nil, errors.New("type assertion error")
	}

	if kid, ok := jwtToken.Header["kid"].(string); ok {
		key, err := rs.publicKeyGetter.getPublicKey(kid)
		if err != nil {
			return nil, err
		}
		return key, nil
	}
	return nil, errors.New("token doesn't have kid header")
}
