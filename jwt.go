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
	rs256SigningKey SigningKeyGetter
	hs256SigningKey SigningKeyGetter
}

func NewMasterValidator(config Config, client *LoginApi) *MasterValidator {
	cks := NewCachedValidationKeysStorage(*client, config.Cache)
	rsa := RSAPublicKey{storage: cks, projectId: config.LoginProjectId}

	return &MasterValidator{
		Config:          config,
		rs256SigningKey: RS256SigningKey{config, rsa},
		hs256SigningKey: HS256SigningKey{config.ShaSecretKey},
	}
}

func (mv MasterValidator) Validate(tokenString string) error {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		signingMethod := token.Method
		switch signingMethod {
		case jwt.SigningMethodRS256:
			return mv.rs256SigningKey.getKey(token)
		case jwt.SigningMethodHS256:
			return mv.hs256SigningKey.getKey(token)
		default:
			return nil, errors.New(signingMethod.Alg() + " algorithm is not supported")
		}
	})

	return err
}
