package login_sdk_go

import (
	"context"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"

	"gitlab.loc/sdk-login/login-sdk-go/contract"
	sdkcontract "gitlab.loc/sdk-login/login-sdk-go/interfaces"
	"gitlab.loc/sdk-login/login-sdk-go/internal/contract"
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
	ValidateWithClaims(ctx context.Context, token string, claims sdkcontract.Claims) (*jwt.Token, error)
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

func (mv *MasterValidator) ValidateWithClaims(ctx context.Context, token string, claims sdkcontract.Claims) (*jwt.Token, error) {
	return mv.validateToken(ctx, token, claims)
}

func (mv *MasterValidator) Validate(ctx context.Context, tokenString string) (*jwt.Token, error) {
	return mv.validateToken(ctx, tokenString, &CustomClaims{})
}

func (mv *MasterValidator) validateToken(ctx context.Context, token string, claims sdkcontract.Claims) (*jwt.Token, error) {
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
		return nil, ErrFailedParseJWT
	}

	switch validationErr.Inner {
	case errSHASecretKeyIsEmpty:
		return mv.validateViaAPI(ctx, parsedToken, token)
	default:
		return nil, err
	}
}

func (mv *MasterValidator) getParserKeyFunction(ctx context.Context) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		signingMethod := token.Method
		switch signingMethod {
		case jwt.SigningMethodRS256:
			return mv.rs256SigningKey.getKey(ctx, token)
		case jwt.SigningMethodHS256:
			return mv.hs256SigningKey.getKey(ctx, token)
		default:
			return nil, fmt.Errorf("%w:%s", ErrNotSupportedAlgorithm, signingMethod.Alg())
		}
	}
}

func (mv *MasterValidator) validateViaAPI(ctx context.Context, parsedToken *jwt.Token, tokenString string) (*jwt.Token, error) {
	err := mv.hs256LoginAPIValidator.Validate(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("mv.hs256LoginAPIValidator.Validate: %w", err)
	}

	if err = validateTokenClaims(parsedToken); err != nil {
		return nil, err
	}

	return parsedToken, nil
}

func validateTokenClaims(parsedToken *jwt.Token) error {
	if err := parsedToken.Claims.Valid(); err != nil {
		return fmt.Errorf("invalid token claims: %w", err)
	}

	return nil
}

func (hs HS256LoginAPIValidator) Validate(ctx context.Context, token string) error {
	return hs.loginAPI.ValidateHS256Token(ctx, token)
}
