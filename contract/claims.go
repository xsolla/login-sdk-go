package contract

import "github.com/dgrijalva/jwt-go"

// Claims interface that must be implemented to use custom claims.
type Claims interface {
	jwt.Claims
	GetProjectID() string
}
