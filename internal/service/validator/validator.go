package validator

import (
	"context"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"

	"gitlab.loc/sdk-login/login-sdk-go/internal/contract"
	"gitlab.loc/sdk-login/login-sdk-go/internal/service/apivalidator"
	"gitlab.loc/sdk-login/login-sdk-go/keys"
)

type Config struct {
	ShaSecretKey string
	Cache        contract.ValidationKeysCache
}

type Validator struct {
	rs256SigningKey        signingKeyGetter
	hs256SigningKey        signingKeyGetter
	hs256LoginApiValidator loginAPIValidator
}

func New(config Config, loginAPI contract.LoginAPI) (*Validator, error) {
	cks := keys.NewCachedValidationKeysStorage(loginAPI, config.Cache)
	rsa := keys.NewRSAPublicKeyGetter(cks)

	return &Validator{
		rs256SigningKey:        keys.NewRS256SigningKeyGetter(rsa),
		hs256SigningKey:        keys.NewHS256SigningKeyGetter(config.ShaSecretKey),
		hs256LoginApiValidator: apivalidator.New(loginAPI),
	}, nil
}

func (mv Validator) Validate(ctx context.Context, token string, claims contract.SDKClaims) (*jwt.Token, error) {
	return mv.validateToken(ctx, token, claims)
}

func (mv Validator) validateToken(ctx context.Context, token string, claims contract.SDKClaims) (*jwt.Token, error) {
	parsedToken, err := jwt.ParseWithClaims(token, claims, mv.getParserKeyFunction(ctx))
	if err == nil {
		if err = validateTokenClaims(parsedToken); err != nil {
			return nil, err
		}

		return parsedToken, nil
	}

	// If there is no secret, validation via API
	// return an error if token is invalid
	validationErr, ok := err.(*jwt.ValidationError)
	if !ok {
		return nil, errors.New("failed parse jwt validation error")
	}

	if errors.Is(validationErr.Inner, keys.ErrSHASecretKeyIsEmpty) {
		if err = mv.hs256LoginApiValidator.Validate(ctx, token); err != nil {
			return nil, err
		}

		if err = validateTokenClaims(parsedToken); err != nil {
			return nil, err
		}

		return parsedToken, nil
	}

	return nil, err
}

func (mv Validator) getParserKeyFunction(ctx context.Context) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		signingMethod := token.Method
		switch signingMethod {
		case jwt.SigningMethodRS256:
			return mv.rs256SigningKey.GetKey(ctx, token)
		case jwt.SigningMethodHS256:
			return mv.hs256SigningKey.GetKey(ctx, token)
		default:
			return nil, fmt.Errorf("algorithm %s is not supported", signingMethod.Alg())
		}
	}
}
