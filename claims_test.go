package login_sdk_go

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCustomClaims(t *testing.T) {
	t.Run("should require jti claim for oauth2 tokens", func(t *testing.T) {
		oauthToken := CustomClaims{ProjectId: "42", Type: "oauth", StandardClaims: jwt.StandardClaims{Id: "123"}}
		err := oauthToken.Valid()
		assert.NoError(t, err)

		anotherToken := CustomClaims{ProjectId: "42", Type: "oauth"}
		err = anotherToken.Valid()
		assert.Error(t, err)
	})

	t.Run("should NOT require jti claim for NOT oauth tokens", func(t *testing.T) {
		oauthToken := CustomClaims{ProjectId: "42", Type: "any", StandardClaims: jwt.StandardClaims{Id: "123"}}
		err := oauthToken.Valid()
		assert.NoError(t, err)

		anotherToken := CustomClaims{ProjectId: "42", Type: "any"}
		err = anotherToken.Valid()
		assert.NoError(t, err)
	})
}
