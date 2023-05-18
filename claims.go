package login_sdk_go

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrJTIClaimIsRequired    = errors.New("jti claim is required")
	ErrXsollaClaimIsRequired = errors.New("xsolla_login_project_id claim is required")
)

type CustomClaims struct {
	ProjectID string   `json:"xsolla_login_project_id,omitempty"`
	Type      string   `json:"type,omitempty"`
	Audience  []string `json:"aud,omitempty"`
	jwt.RegisteredClaims
}

func (c CustomClaims) Valid() error {
	if c.ID == "" && c.Type == "oauth" {
		return ErrJTIClaimIsRequired
	}

	if c.ProjectID == "" {
		return ErrXsollaClaimIsRequired
	}

	return nil
}

func (c CustomClaims) GetProjectID() string {
	return c.ProjectID
}
