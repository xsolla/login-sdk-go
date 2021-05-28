package login_sdk_go

import (
	"crypto/rsa"
	"errors"
	"github.com/dgrijalva/jwt-go"
)

type Validator interface {
	Validate(jwt string) error
}

type MasterValidator struct {
	Options
}

func (mv MasterValidator) Validate(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		signingMethod := token.Method
		switch signingMethod {
		case jwt.SigningMethodRS256:
			loginApi := LoginApi{baseUrl: mv.LoginApiUrl}
			keysResp, _ := loginApi.GetProjectKeysForLoginProject(mv.LoginProjectId)

			// todo: refactor
			var pubKey RSAKeyResponse
			if kid, ok := token.Header["kid"]; ok {
				if len(keysResp) > 1 {
					for i := range keysResp {
						if keysResp[i].Kid == kid {
							pubKey = keysResp[i]
							break
						}
					}
				} else {
					pubKey = keysResp[0]
				}
			}

			return &rsa.PublicKey{
				N: fromBase16(pubKey.Modulus),
				E: int(fromBase16(pubKey.Exponent).Int64()),
			}, nil

		case jwt.SigningMethodHS256:
			return []byte(mv.ShaSecretKey), nil
		default:
			return nil, errors.New("not supported algorithm")
		}
	})

	return token, err
}
