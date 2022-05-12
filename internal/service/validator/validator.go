package validator

import (
	"context"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"

	sdkcontract "gitlab.loc/sdk-login/login-sdk-go/contract"
	"gitlab.loc/sdk-login/login-sdk-go/internal/contract"
	"gitlab.loc/sdk-login/login-sdk-go/internal/service/apivalidator"
	"gitlab.loc/sdk-login/login-sdk-go/keys"
)

var (
	ErrFailedParseJWT        = errors.New("failed parse jwt validation error")
	ErrNotSupportedAlgorithm = errors.New("algorithm is not supported")
)

type Config struct {
	ShaSecretKey string
	Cache        sdkcontract.ValidationKeysCache
}

type Validator struct {
	rs256SigningKey        signingKeyGetter
	hs256SigningKey        signingKeyGetter
	hs256LoginAPIValidator loginAPIValidator
}

func New(config Config, loginAPI contract.LoginAPI) (*Validator, error) {
	cks := keys.NewCachedValidationKeysStorage(loginAPI, config.Cache)
	rsa := keys.NewRSAPublicKeyGetter(cks)

	return &Validator{
		rs256SigningKey:        keys.NewRS256SigningKeyGetter(rsa),
		hs256SigningKey:        keys.NewHS256SigningKeyGetter(config.ShaSecretKey),
		hs256LoginAPIValidator: apivalidator.New(loginAPI),
	}, nil
}

func (v *Validator) Validate(ctx context.Context, token string, claims sdkcontract.Claims) (*jwt.Token, error) {
	return v.validateToken(ctx, token, claims)
}

func (v *Validator) validateToken(ctx context.Context, token string, claims sdkcontract.Claims) (*jwt.Token, error) {
	parsedToken, err := jwt.ParseWithClaims(token, claims, v.getParserKeyFunction(ctx))
	if err == nil {
		if err = validateTokenClaims(parsedToken); err != nil {
			return nil, err
		}

		return parsedToken, nil
	}

	// If there is no secret, validation via API
	// return an error if token is invalid
	// If there is no secret, validation via API
	// return an error if token is invalid
	validationErr, ok := err.(*jwt.ValidationError)
	if !ok {
		return nil, ErrFailedParseJWT
	}

	switch validationErr.Inner {
	case keys.ErrSHASecretKeyIsEmpty:
		return v.validateViaAPI(ctx, parsedToken, token)
	default:
		return nil, err
	}
}

func (v *Validator) getParserKeyFunction(ctx context.Context) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		signingMethod := token.Method
		switch signingMethod {
		case jwt.SigningMethodRS256:
			return v.rs256SigningKey.GetKey(ctx, token)
		case jwt.SigningMethodHS256:
			return v.hs256SigningKey.GetKey(ctx, token)
		default:
			return nil, fmt.Errorf("%w:%s", ErrNotSupportedAlgorithm, signingMethod.Alg())
		}
	}
}

func (v *Validator) validateViaAPI(ctx context.Context, parsedToken *jwt.Token, tokenString string) (*jwt.Token, error) {
	err := v.hs256LoginAPIValidator.Validate(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("v.hs256LoginAPIValidator.Validate: %w", err)
	}

	if err = validateTokenClaims(parsedToken); err != nil {
		return nil, err
	}

	return parsedToken, nil
}
