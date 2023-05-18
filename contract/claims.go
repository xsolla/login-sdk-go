package contract

import "github.com/golang-jwt/jwt/v5"

// Claims interface that must be implemented to use custom claims.
type Claims interface {
	GetSubject() (string, error)
	GetIssuer() (string, error)
	GetNotBefore() (*jwt.NumericDate, error)
	GetIssuedAt() (*jwt.NumericDate, error)
	GetExpirationTime() (*jwt.NumericDate, error)
	GetAudience() (jwt.ClaimStrings, error)
	GetProjectID() string
}
