package validator

import (
	"context"

	"github.com/dgrijalva/jwt-go"
)

type loginAPIValidator interface {
	Validate(ctx context.Context, jwt string) error
}

type signingKeyGetter interface {
	GetKey(ctx context.Context, token *jwt.Token) (interface{}, error)
}
