package login_sdk_go

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

type TokenString string

var (
	ErrJTIClaimIsRequired    = errors.New("jti claim is required")
	ErrXsollaClaimIsRequired = errors.New("xsolla_login_project_id claim is required")
)

type CustomClaims struct {
	ProjectID string   `json:"xsolla_login_project_id,omitempty"`
	Audience  []string `json:"aud,omitempty"`
	Type      string   `json:"type:omitempty"`
	jwt.StandardClaims
}

func (c CustomClaims) Valid() error {
	vErr := new(jwt.ValidationError)

	if c.Id == "" && c.Type == "oauth" {
		vErr.Inner = ErrJTIClaimIsRequired
		vErr.Errors |= jwt.ValidationErrorId
	}

	if c.ProjectID == "" {
		vErr.Inner = ErrXsollaClaimIsRequired
		vErr.Errors |= jwt.ValidationErrorClaimsInvalid
	}

	if vErr.Inner != nil {
		return vErr
	}

	return c.StandardClaims.Valid()
}

func (c CustomClaims) GetProjectID() string {
	return c.ProjectId
}
