package login_sdk_go

import (
	"github.com/dgrijalva/jwt-go"
	"time"

	"gitlab.loc/sdk-login/login-sdk-go/cache"
	"gitlab.loc/sdk-login/login-sdk-go/infrastructure"
	"gitlab.loc/sdk-login/login-sdk-go/model"
)

const (
	defaultLoginApiUrl = "https://login.xsolla.com"
)

type Config struct {
	IgnoreSslErrors        bool
	ShaSecretKey           string
	LoginApiUrl            string
	LoginProjectId         string
	LoginClientId          int
	LoginClientSecret      string
	SessionApiHost         string
	SessionApiPort         int
	SkipSessionValidation  bool
	IsMultipleProjectsMode bool
	Cache                  cache.ValidationKeysCache
}

type ConfigOption func(*Config)

type LoginSdk struct {
	config    Config
	validator ValidatorWithParser
	refresher Refresher
}

func New(config Config) (*LoginSdk, error) {
	config.fillDefaults()

	loginApi := infrastructure.NewHttpLoginApi(config.LoginApiUrl, config.IgnoreSslErrors)

	mv, err := NewMasterValidator(config, &loginApi)

	if err != nil {
		return nil, err
	}

	l := &LoginSdk{
		config:    config,
		validator: mv,
		refresher: NewTokenRefresher(&loginApi, config.LoginClientId, config.LoginClientSecret),
	}

	return l, nil
}

func (c *Config) fillDefaults() {
	if c.LoginApiUrl == "" {
		c.LoginApiUrl = defaultLoginApiUrl
	}

	if c.Cache == nil {
		c.Cache = cache.NewDefaultCache(1 * time.Hour)
	}
}

func (sdk *LoginSdk) Validate(tokenString string) (*jwt.Token, *WrappedError) {
	parsedToken, err := sdk.validator.Validate(tokenString)
	return parsedToken, WrapError(err)
}

func (sdk LoginSdk) Refresh(refreshToken string) (*model.LoginToken, error) {
	return sdk.refresher.Refresh(refreshToken)
}
