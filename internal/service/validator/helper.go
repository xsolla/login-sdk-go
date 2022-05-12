package validator

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

func validateTokenClaims(parsedToken *jwt.Token) error {
	if err := parsedToken.Claims.Valid(); err != nil {
		return fmt.Errorf("invalid token claims: %w", err)
	}

	return nil
}
