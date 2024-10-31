package utils

import (
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(claims jwt.Claims, secretKey string) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(secretKey))

	return tokenString, err
}
