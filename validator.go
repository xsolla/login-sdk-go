package login_sdk_go

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"gitlab.loc/sdk-login/login-sdk-go/interfaces"
)

type Validator interface {
	Validate(jwt string) error
}

type ValidatorWithParser interface {
	Validate(jwt string) (*jwt.Token, error)
}

type HS256LoginApiValidator struct {
	loginApi *interfaces.LoginApi
}

type MasterValidator struct {
	Config
	rs256SigningKey        SigningKeyGetter
	hs256SigningKey        SigningKeyGetter
	hs256LoginApiValidator Validator
}

func NewMasterValidator(config Config, client *interfaces.LoginApi) (*MasterValidator, error) {
	cks := NewCachedValidationKeysStorage(*client, config.Cache)
	rsa := RSAPublicKeyGetter{storage: cks}

	hs256LoginApiValidator := &HS256LoginApiValidator{client}

	return &MasterValidator{
		Config:                 config,
		rs256SigningKey:        RS256SigningKeyGetter{config, rsa},
		hs256SigningKey:        HS256SigningKeyGetter{config.ShaSecretKey},
		hs256LoginApiValidator: hs256LoginApiValidator,
	}, nil
}

func (mv MasterValidator) Validate(tokenString string) (*jwt.Token, error) {
	parsedToken, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
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

	if err != nil {
		// Если секрета нет, валидируем через апишку
		// и вернем ошибку, если токен не валиден
		validationErr, ok := err.(*jwt.ValidationError)
		if !ok {
			return nil, errors.New("failed parse jwt validation error")
		}
		if validationErr.Inner == errSHASecretKeyIsEmpty {
			err = mv.hs256LoginApiValidator.Validate(tokenString)
			if err != nil {
				return nil, err
			}
			err = parsedToken.Claims.Valid()
			if err != nil {
				return nil, err
			}
			return parsedToken, nil
		}
		return nil, err
	}
	// подпись валидная, проверяем только истек токен или нет
	err = parsedToken.Claims.Valid()
	if err != nil {
		return nil, err
	} else {
		return parsedToken, nil
	}

}

func (hs HS256LoginApiValidator) Validate(token string) error {
	l := *hs.loginApi
	return l.ValidateHS256Token(token)
}
