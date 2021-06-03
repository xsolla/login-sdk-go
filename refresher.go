package login_sdk_go

import (
	"gitlab.loc/sdk-login/login-sdk-go/interfaces"
	"gitlab.loc/sdk-login/login-sdk-go/model"
)

type Refresher interface {
	Refresh(s string) (*model.LoginToken, error)
}

type tokenRefresher struct {
	loginApi     *interfaces.LoginApi
	clientId     int
	clientSecret string
}

func NewTokenRefresher(loginApi *interfaces.LoginApi, clientId int, clientSecret string) *tokenRefresher {
	return &tokenRefresher{
		loginApi:     loginApi,
		clientId:     clientId,
		clientSecret: clientSecret,
	}
}

func (r tokenRefresher) Refresh(token string) (*model.LoginToken, error) {
	l := *r.loginApi
	return l.RefreshToken(token, r.clientId, r.clientSecret)
}
