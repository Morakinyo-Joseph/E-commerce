package models

import "github.com/dgrijalva/jwt-go"

// JwtClaims represents the structure of the JWT claims.
type JwtClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`

	jwt.StandardClaims
}
