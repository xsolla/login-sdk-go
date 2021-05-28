package login_sdk_go

import (
	"github.com/dgrijalva/jwt-go"
)

const (
	defaultLoginApiUrl = "https://login.xsolla.com"
	defaultIssuer      = "https://login.xsolla.com"
)

type Options struct {
	ShaSecretKey      string
	LoginApiUrl       string
	LoginProjectId    string
	Issuer            string
	LoginClientId     int
	LoginClientSecret string
}

type loginSdk struct {
	Options
}

func New(options Options) *loginSdk {
	options.fillDefaults()

	l := &loginSdk{
		options,
	}

	return l
}

func (o *Options) fillDefaults() {
	if o.LoginApiUrl == "" {
		o.LoginApiUrl = defaultLoginApiUrl
	}

	if o.Issuer == "" {
		o.Issuer = defaultIssuer
	}
}

func (sdk loginSdk) Validate(tokenString string) (*jwt.Token, error) {
	token, err := MasterValidator{sdk.Options}.Validate(tokenString)
	return token, err
}

func (sdk loginSdk) Refresh(refreshToken string) (*LoginTokenResponse, error) {
	loginApi := LoginApi{baseUrl: sdk.LoginApiUrl}
	response, err := loginApi.RefreshToken(refreshToken, sdk.LoginClientId, sdk.LoginClientSecret)
	return &response, err
}
