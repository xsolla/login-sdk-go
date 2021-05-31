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
			if kid, ok := token.Header["kid"].(string); ok {
				client := NewHttpLoginApi(mv.LoginApiUrl)
				cks := NewCachedValidationKeysStorage(client, mv.Cache)
				rsa := RSASigningKey{storage: cks, projectId: mv.LoginProjectId}
				key, err := rsa.getSigningKey(kid)
				if err != nil {
					return nil, err
				}
				return key, nil
			}
			return nil, errors.New("token doesn't have kid header")
		case jwt.SigningMethodHS256:
			return []byte(mv.ShaSecretKey), nil
		default:
			return nil, errors.New(signingMethod.Alg() + " algorithm is not supported")
		}
	})

	return token, err
}
