package login_sdk_go

import (
	"github.com/golang-jwt/jwt/v5"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomClaims(t *testing.T) {
	t.Run("should require jti claim for oauth2 tokens", func(t *testing.T) {
		oauthToken := CustomClaims{ProjectID: "42", Type: "oauth", RegisteredClaims: jwt.RegisteredClaims{ID: "123"}}
		err := oauthToken.Valid()
		assert.NoError(t, err)

		anotherToken := CustomClaims{ProjectID: "42", Type: "oauth"}
		err = anotherToken.Valid()
		assert.Error(t, err)
	})

	t.Run("should NOT require jti claim for NOT oauth tokens", func(t *testing.T) {
		oauthToken := CustomClaims{ProjectID: "42", Type: "any", RegisteredClaims: jwt.RegisteredClaims{ID: "123"}}
		err := oauthToken.Valid()
		assert.NoError(t, err)

		anotherToken := CustomClaims{ProjectID: "42", Type: "any"}
		err = anotherToken.Valid()
		assert.NoError(t, err)
	})
}
