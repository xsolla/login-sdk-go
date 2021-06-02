package login_sdk_go

import (
	"time"
)

const (
	defaultLoginApiUrl = "https://login.xsolla.com"
)

type Config struct {
	IgnoreSslErrors   bool
	ShaSecretKey      string
	LoginApiUrl       string
	LoginProjectId    string
	LoginClientId     int
	LoginClientSecret string
	Cache             Cache
}

type ConfigOption func(*Config)

type loginSdk struct {
	config    Config
	validator Validator
	refresher Refresher
	loginApi  *LoginApi
}

func New(config Config) *loginSdk {
	config.fillDefaults()

	loginApi := NewHttpLoginApi(config.LoginApiUrl, config.IgnoreSslErrors)

	l := &loginSdk{
		config:    config,
		validator: NewMasterValidator(config, &loginApi),
		refresher: NewTokenRefresher(&loginApi, config.LoginClientId, config.LoginClientSecret),
	}

	return l
}

func (c *Config) fillDefaults() {
	if c.LoginApiUrl == "" {
		c.LoginApiUrl = defaultLoginApiUrl
	}

	if c.Cache == nil {
		c.Cache = NewDefaultCache(1 * time.Hour)
	}
}

func (sdk *loginSdk) Validate(tokenString string) *WrappedError {
	err := sdk.validator.Validate(tokenString)
	return WrapError(err)
}

func (sdk loginSdk) Refresh(refreshToken string) (*LoginToken, error) {
	return sdk.refresher.Refresh(refreshToken)
}
