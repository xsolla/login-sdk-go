package login_sdk_go

type Refresher interface {
	Refresh(s string) (*LoginToken, error)
}

type tokenRefresher struct {
	loginApi     *LoginApi
	clientId     int
	clientSecret string
}

func NewTokenRefresher(loginApi *LoginApi, clientId int, clientSecret string) *tokenRefresher {
	return &tokenRefresher{
		loginApi:     loginApi,
		clientId:     clientId,
		clientSecret: clientSecret,
	}
}

func (r tokenRefresher) Refresh(token string) (*LoginToken, error) {
	l := *r.loginApi
	return l.RefreshToken(token, r.clientId, r.clientSecret)
}
