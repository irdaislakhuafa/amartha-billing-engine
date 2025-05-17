package entity

import "github.com/golang-jwt/jwt/v5"

type (
	JWTClaims struct {
		UID int64
		jwt.RegisteredClaims
	}
)
