package login_sdk_go

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"

	"gitlab.loc/sdk-login/login-sdk-go/cache"
	"gitlab.loc/sdk-login/login-sdk-go/contract"
	"gitlab.loc/sdk-login/login-sdk-go/internal/adapter/login"
	vl "gitlab.loc/sdk-login/login-sdk-go/internal/service/validator"
)

const (
	defaultLoginAPIURL = "https://login.xsolla.com"
	keyTTL             = 10 * time.Minute
	defaultAPITimeout  = 5 * time.Second
)

type Config struct {
	IgnoreSslErrors bool
	ShaSecretKey    string
	LoginAPIURL     string
	Cache           contract.ValidationKeysCache
	APITimeout      time.Duration
}

type ConfigOption func(*Config)

type validator interface {
	Validate(ctx context.Context, jwt string, claims contract.Claims) (*jwt.Token, error)
}

type LoginSdk struct {
	config    Config
	validator validator
}

func New(config Config) (*LoginSdk, error) {
	config.fillDefaults()

	loginApi := login.NewAdapter(config.LoginAPIURL, config.IgnoreSslErrors, config.APITimeout)

	validator, err := vl.New(vl.Config{
		ShaSecretKey: config.ShaSecretKey,
		Cache:        config.Cache,
	}, loginApi)
	if err != nil {
		return nil, err
	}

	l := &LoginSdk{
		config:    config,
		validator: validator,
	}

	return l, nil
}

func (c *Config) fillDefaults() {
	if c.LoginAPIURL == "" {
		c.LoginAPIURL = defaultLoginAPIURL
	}

	if c.Cache == nil {
		c.Cache = cache.NewDefaultCache(keyTTL)
	}
	if c.APITimeout == 0 {
		c.APITimeout = defaultAPITimeout
	}
}

func (sdk *LoginSdk) ValidateWithContext(ctx context.Context, token string) (*jwt.Token, *WrappedError) {
	parsedToken, err := sdk.validator.Validate(ctx, token, &CustomClaims{})

	return parsedToken, WrapError(err)
}

func (sdk *LoginSdk) Validate(token string) (*jwt.Token, *WrappedError) {
	return sdk.ValidateWithContext(context.Background(), token)
}

func (sdk *LoginSdk) ValidateWithClaimsAndContext(
	ctx context.Context,
	token string,
	claims contract.Claims,
) (*jwt.Token, *WrappedError) {
	parsedToken, err := sdk.validator.Validate(ctx, token, claims)
	if err != nil {
		return nil, WrapError(err)
	}

	return parsedToken, nil
}

func (sdk *LoginSdk) ValidateWithClaims(token string, claims contract.Claims) (*jwt.Token, *WrappedError) {
	return sdk.ValidateWithClaimsAndContext(context.Background(), token, claims)
}
