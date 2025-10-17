package helper

import (
	"WedhaWS/utsgosmt5/alumni-tracer/app/models/domain"
	"WedhaWS/utsgosmt5/alumni-tracer/app/models/web"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key-min-32-characters-long")

func GenerateToken(user domain.Users) (string, error) {
	claims := web.JWTClaims{
		UserID:   user.Id,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*web.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &web.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*web.JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrInvalidKey
}
