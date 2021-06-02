package login_sdk_go

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateMalformedToken(t *testing.T) {
	token := ""
	loginSgk := New(Config{})
	err := loginSgk.Validate(token)

	require.False(t, err.Valid())
}

func TestValidateHmacToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	loginSgk := New(Config{
		ShaSecretKey: "your-256-bit-secret",
	})
	err := loginSgk.Validate(token)

	require.True(t, err.Valid())

	if !err.Valid() {
		t.Fatal(err)
	}
}

func TestValidateExpiredRsaToken(t *testing.T) {
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6InNnRnk0NjRrVk5YVFo2YmVYM0tFT2kyam1yWnA4bUQiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOltdLCJlbWFpbCI6ImxvZ2luLXNka0B0ZXN0LmNvbSIsImV4cCI6MTYyMjQ2MTU5MSwiZ3JvdXBzIjpbXSwiaWF0IjoxNjIyNDU3OTkwLCJpc19tYXN0ZXIiOnRydWUsImlzcyI6Imh0dHBzOi8vbG9naW4ueHNvbGxhLmNvbSIsImp0aSI6ImZmZjU3MTNmLTY2YWItNGVmNy04OGVlLTIyYWFmNmU0YTUzNyIsInByb21vX2VtYWlsX2FncmVlbWVudCI6dHJ1ZSwic2NwIjpbIm9mZmxpbmUiXSwic3ViIjoiOTNlMWE5YTMtOGRkZC00OTMxLWE4ZmQtMjA5MGY2M2VkZmI4IiwidHlwZSI6Inhzb2xsYV9sb2dpbiIsInVzZXJuYW1lIjoibG9naW4tc2RrLXRlc3QiLCJ4c29sbGFfbG9naW5fYWNjZXNzX2tleSI6ImZPOGNPc2o4TjdHSWw3MldscVNsYndnSmk0ZGRBOHVQbFd2ZDYxTk1ydnMiLCJ4c29sbGFfbG9naW5fcHJvamVjdF9pZCI6IjQwZGIyZWE0LTVkNDItMTFlNi1hM2ZmLTAwNTA1NmEwZTA0YSJ9.PaOD3zvCVxqhpTUb1qD9zhnbtumRVoBqpLTCyvhS0cOF-zl2hpCnZO_brJgR0eTNyaAcYdR1Q87ZxdUvs_ILsXAdPLNoazK22X4tO17VV6HsqyIx47_6KHeNU42Z8PXY59ZDXTSWU7m9J1-8vRnuGYN10QNmr7kqPR21i-xsr8D0pV2sn23GWgvlZr3m2-zZQTRBu_xE2IMtwyaTDCUweY7NWvGPB8O3E_CBDXOxH6DD6j42mLZ7XabjWTcd8zwQNNA9kmmAEN3D3LMwuKeU2xHhM6eRT6--xTccr8jTeuQDetWekdSQRZxGWOJlq_iyO_dne7snwOCSQ9MP1GbFGg"
	loginSgk := New(Config{
		LoginProjectId:  "40db2ea4-5d42-11e6-a3ff-005056a0e04a",
		IgnoreSslErrors: true,
	})
	err := loginSgk.Validate(token)

	require.False(t, err.Valid())
	require.True(t, err.IsExpired())
}

type testRefresher struct {
	result *LoginToken
}

func (r testRefresher) Refresh(t string) (*LoginToken, error) {
	return r.result, nil
}

func TestRefreshToken(t *testing.T) {
	refreshToken := "YgAhlbpA6_y63D_OCOygdloo_ci87h8b5tF93oKasqM.w2PFSamrNadObdcZLCIneA0XHe8zae1EWsURh_6-eIg"

	expected := &LoginToken{
		RefreshToken: "123",
		AccessToken:  "123",
		ExpiresIn:    123,
	}

	loginSgk := loginSdk{
		refresher: testRefresher{result: expected},
	}

	res, err := loginSgk.Refresh(refreshToken)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, expected, res)
}
