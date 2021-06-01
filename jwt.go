package login_sdk_go

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

type Validator interface {
	Validate(jwt string) error
}

type SigningKeyGetter interface {
	getKey(token interface{}) (interface{}, error)
}

type MasterValidator struct {
	Config
	rs256Validator SigningKeyGetter
	hs256Validator SigningKeyGetter
}

type HS256Validator struct {
	key string
}

type RS256Validator struct {
	Config
}

func NewMasterValidator(config Config) *MasterValidator {
	return &MasterValidator{
		Config:         config,
		rs256Validator: RS256Validator{config},
		hs256Validator: HS256Validator{config.ShaSecretKey},
	}
}

func (mv MasterValidator) Validate(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		signingMethod := token.Method
		switch signingMethod {
		case jwt.SigningMethodRS256:
			return mv.rs256Validator.getKey(token)
		case jwt.SigningMethodHS256:
			return mv.hs256Validator.getKey(token)
		default:
			return nil, errors.New(signingMethod.Alg() + " algorithm is not supported")
		}
	})

	return token, err
}

func (hs HS256Validator) getKey(interface{}) (interface{}, error) {
	if hs.key == "" {
		return nil, errors.New("ShaSecretKey config is required for sync validation")
	}

	return []byte(hs.key), nil
}

func (rs RS256Validator) getKey(token interface{}) (interface{}, error) {
	jwtToken, ok := token.(*jwt.Token)
	if ok == false {
		return nil, errors.New("type assertion error")
	}

	if kid, ok := jwtToken.Header["kid"].(string); ok {
		client := NewHttpLoginApi(rs.LoginApiUrl)
		cks := NewCachedValidationKeysStorage(client, rs.Cache)
		rsa := RSAPublicKey{storage: cks, projectId: rs.LoginProjectId}
		key, err := rsa.getPublicKey(kid)
		if err != nil {
			return nil, err
		}
		return key, nil
	}
	return nil, errors.New("token doesn't have kid header")
}
