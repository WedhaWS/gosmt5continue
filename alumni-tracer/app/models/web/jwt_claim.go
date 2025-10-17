package web

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID int64 `json:"user_id"`
	Username string `json:"username"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}
