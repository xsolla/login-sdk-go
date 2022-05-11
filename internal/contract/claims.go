package contract

import "github.com/dgrijalva/jwt-go"

// SDKClaims interface that must be implemented to use custom claims.
type SDKClaims interface {
	jwt.Claims
	GetProjectID() string
}
