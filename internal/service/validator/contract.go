package validator

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type loginAPIValidator interface {
	Validate(ctx context.Context, jwt string) error
}

type signingKeyGetter interface {
	GetKey(ctx context.Context, token *jwt.Token) (interface{}, error)
}
