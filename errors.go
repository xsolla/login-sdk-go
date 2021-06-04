package login_sdk_go

import (
	"github.com/dgrijalva/jwt-go"
)

type WrappedError struct {
	Inner error
}

func WrapError(err error) *WrappedError {
	return &WrappedError{
		err,
	}
}

func (we *WrappedError) Error() string {
	return we.Inner.Error()
}

func (we *WrappedError) Valid() bool {
	return we.Inner == nil
}

func (we *WrappedError) IsExpired() bool {
	ve, ok := we.Inner.(*jwt.ValidationError)
	if ok {
		return jwt.ValidationErrorExpired == ve.Errors
	}

	return false
}
