package login_sdk_go

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	defaultLoginApiUrl = "https://login.xsolla.com"
	defaultIssuer      = "https://login.xsolla.com"
)

type Config struct {
	ShaSecretKey      string
	LoginApiUrl       string
	LoginProjectId    string
	Issuer            string
	LoginClientId     int
	LoginClientSecret string
	Cache             Cache
}

type ConfigOption func(*Config)

type loginSdk struct {
	Config
}

func New(config Config) *loginSdk {
	config.fillDefaults()

	l := &loginSdk{
		config,
	}

	return l
}

func (c *Config) fillDefaults() {
	if c.LoginApiUrl == "" {
		c.LoginApiUrl = defaultLoginApiUrl
	}

	if c.Issuer == "" {
		c.Issuer = defaultIssuer
	}

	if c.Cache == nil {
		c.Cache = NewDefaultCache(1 * time.Minute)
	}
}

func (sdk loginSdk) Validate(tokenString string) (*jwt.Token, error) {
	token, err := MasterValidator{sdk.Config}.Validate(tokenString)
	return token, err
}

func (sdk loginSdk) Refresh(refreshToken string) (*LoginToken, error) {
	loginApi := httpLoginApi{baseUrl: sdk.LoginApiUrl}
	response, err := loginApi.RefreshToken(refreshToken, sdk.LoginClientId, sdk.LoginClientSecret)
	return &response, err
}
