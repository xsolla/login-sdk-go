package login_sdk_go

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

type Validator interface {
	Validate(jwt string) error
}

type MasterValidator struct {
	Config
}

func (mv MasterValidator) Validate(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		signingMethod := token.Method
		switch signingMethod {
		case jwt.SigningMethodRS256:
			return getRS256SigningKey(token, mv.Config)
		case jwt.SigningMethodHS256:
			return getHS256SigningKey(mv.Config)
		default:
			return nil, errors.New(signingMethod.Alg() + " algorithm is not supported")
		}
	})

	return token, err
}

func getHS256SigningKey(c Config) (interface{}, error) {
	if c.ShaSecretKey == "" {
		return nil, errors.New("ShaSecretKey config is required for sync validation")
	}

	return []byte(c.ShaSecretKey), nil
}

func getRS256SigningKey(token *jwt.Token, c Config) (interface{}, error) {
	if kid, ok := token.Header["kid"].(string); ok {
		client := NewHttpLoginApi(c.LoginApiUrl)
		cks := NewCachedValidationKeysStorage(client, c.Cache)
		rsa := RSAPublicKey{storage: cks, projectId: c.LoginProjectId}
		key, err := rsa.getPublicKey(kid)
		if err != nil {
			return nil, err
		}
		return key, nil
	}
	return nil, errors.New("token doesn't have kid header")
}
