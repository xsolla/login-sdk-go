package apivalidator

import (
	"context"
)

type HS256LoginApiValidator struct {
	loginAPI loginAPI
}

func New(loginAPI loginAPI) HS256LoginApiValidator {
	return HS256LoginApiValidator{
		loginAPI: loginAPI,
	}
}

func (hs HS256LoginApiValidator) Validate(ctx context.Context, token string) error {
	return hs.loginAPI.ValidateHS256Token(ctx, token)
}
