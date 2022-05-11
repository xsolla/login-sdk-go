package apivalidator

import "context"

type loginAPI interface {
	ValidateHS256Token(ctx context.Context, token string) error
}
