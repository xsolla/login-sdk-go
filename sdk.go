package login_sdk_go

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"

	"gitlab.loc/sdk-login/login-sdk-go/cache"
	"gitlab.loc/sdk-login/login-sdk-go/infrastructure"
)

const (
	defaultLoginAPIURL = "https://login.xsolla.com"
	keyTTL             = 10 * time.Minute
)

type Config struct {
	IgnoreSslErrors bool
	ShaSecretKey    string
	LoginAPIURL     string
	Cache           cache.ValidationKeysCache
}

type ConfigOption func(*Config)

type LoginSdk struct {
	config    Config
	validator ValidatorWithParser
}

func New(config Config) (*LoginSdk, error) {
	config.fillDefaults()

	loginAPI := infrastructure.NewHttpLoginAPI(config.LoginAPIURL, config.IgnoreSslErrors)

	validator, err := NewMasterValidator(config, loginAPI)
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
}

func (sdk *LoginSdk) ValidateWithContext(ctx context.Context, tokenString string) (*jwt.Token, *WrappedError) {
	parsedToken, err := sdk.validator.Validate(ctx, tokenString)

	return parsedToken, WrapError(err)
}

func (sdk *LoginSdk) Validate(tokenString string) (*jwt.Token, *WrappedError) {
	return sdk.ValidateWithContext(context.Background(), tokenString)
}
