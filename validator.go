package login_sdk_go

import (
	"context"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"

	"gitlab.loc/sdk-login/login-sdk-go/interfaces"
)

var (
	ErrFailedParseJWT        = errors.New("failed parse jwt validation error")
	ErrNotSupportedAlgorithm = errors.New("algorithm is not supported")
)

type Validator interface {
	Validate(ctx context.Context, jwt string) error
}

type ValidatorWithParser interface {
	Validate(ctx context.Context, jwt string) (*jwt.Token, error)
}

type HS256LoginAPIValidator struct {
	loginAPI interfaces.LoginAPI
}

type MasterValidator struct {
	Config
	rs256SigningKey        SigningKeyGetter
	hs256SigningKey        SigningKeyGetter
	hs256LoginAPIValidator Validator
}

func NewMasterValidator(config Config, client interfaces.LoginAPI) (*MasterValidator, error) {
	cks := NewCachedValidationKeysStorage(client, config.Cache)
	rsa := RSAPublicKeyGetter{storage: cks}

	hs256LoginAPIValidator := &HS256LoginAPIValidator{client}

	return &MasterValidator{
		Config:                 config,
		rs256SigningKey:        RS256SigningKeyGetter{config, rsa},
		hs256SigningKey:        HS256SigningKeyGetter{config.ShaSecretKey},
		hs256LoginAPIValidator: hs256LoginAPIValidator,
	}, nil
}

func (mv MasterValidator) Validate(ctx context.Context, tokenString string) (*jwt.Token, error) {
	parsedToken, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		signingMethod := token.Method
		switch signingMethod {
		case jwt.SigningMethodRS256:
			return mv.rs256SigningKey.getKey(ctx, token)
		case jwt.SigningMethodHS256:
			return mv.hs256SigningKey.getKey(ctx, token)
		default:
			return nil, fmt.Errorf("%w:%s", ErrNotSupportedAlgorithm, signingMethod.Alg())
		}
	})
	if err != nil {
		// If there is no secret, validation via API
		// return an error if token is invalid
		validationErr, ok := err.(*jwt.ValidationError)
		if !ok {
			return nil, ErrFailedParseJWT
		}
		if errors.Is(validationErr.Inner, errSHASecretKeyIsEmpty) {
			err = mv.hs256LoginAPIValidator.Validate(ctx, tokenString)
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
	// sign is valid, checkin token expiration date
	err = parsedToken.Claims.Valid()
	if err != nil {
		return nil, err
	} else {
		return parsedToken, nil
	}
}

func (hs HS256LoginAPIValidator) Validate(ctx context.Context, token string) error {
	l := hs.loginAPI

	return l.ValidateHS256Token(ctx, token)
}
