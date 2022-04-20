package login_sdk_go

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

type TokenString string

type CustomClaims struct {
	ProjectID string   `json:"xsolla_login_project_id,omitempty"`
	Audience  []string `json:"aud,omitempty"`
	Type      string   `json:"type:omitempty"`
	jwt.StandardClaims
}

func (c CustomClaims) Valid() error {
	vErr := new(jwt.ValidationError)

	if c.Id == "" && c.Type == "oauth" {
		vErr.Inner = fmt.Errorf("jti claim is required")
		vErr.Errors |= jwt.ValidationErrorId
	}

	if c.ProjectID == "" {
		vErr.Inner = fmt.Errorf("xsolla_login_project_id claim is required")
		vErr.Errors |= jwt.ValidationErrorClaimsInvalid
	}

	if vErr.Inner != nil {
		return vErr
	}

	return c.StandardClaims.Valid()
}
