package login_sdk_go

import (
	"crypto/rsa"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"log"
	"math/big"
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

func fromBase16(base16 string) *big.Int {
	i, ok := new(big.Int).SetString(base16, 16)
	if !ok {
		log.Fatal("bad number: " + base16)
	}
	return i
}

func (sdk loginSdk) Validate(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		signingMethod := token.Method
		switch signingMethod {
		case jwt.SigningMethodRS256:
			loginApi := LoginApi{baseUrl: sdk.LoginApiUrl}
			keysResp, _ := loginApi.GetProjectKeysForLoginProject(sdk.LoginProjectId)
			pubKey := keysResp[0]

			return &rsa.PublicKey{
				N: fromBase16(pubKey.Modulus),
				E: int(fromBase16(pubKey.Exponent).Int64()),
			}, nil

		case jwt.SigningMethodHS256:
			return []byte(sdk.ShaSecretKey), nil
		default:
			return nil, errors.New("not supported algorithm")
		}
	})

	return token, err
}

func (sdk loginSdk) Refresh(refreshToken string) (*LoginToken, error) {
	loginApi := LoginApi{baseUrl: sdk.LoginApiUrl}
	response, err := loginApi.RefreshToken(refreshToken)
	return &response, err
}
