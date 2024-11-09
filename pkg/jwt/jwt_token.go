package jwt_token

import (
	"delivery-backend/internal/setting"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

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

type RoleClaims struct {
	ID      uint   `json:"id"`
	Account string `json:"account"`
	jwt.RegisteredClaims
}

func GetAccessToken(issuer string, ID uint, account string, expires_minute int) string {
	expires := time.Now().Add(time.Duration(expires_minute) * time.Minute)

	claims := RoleClaims{
		ID,
		account,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
			Issuer:    issuer,
			Subject:   "access_token",
		},
	}

	tks, err := GenerateToken(claims, setting.AppSetting.JWTSecretKey)
	if err != nil {
		// 这里不应当出错
		log.Fatal(err)
	}
	return tks
}

func GetRefreshToken(issuer string, ID uint, account string, expires_minute int) string {
	expires := time.Now().Add(time.Duration(expires_minute) * time.Minute)

	claims := RoleClaims{
		ID,
		account,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
			Issuer:    issuer,
			Subject:   "refresh_token",
		},
	}

	tks, err := GenerateToken(claims, setting.AppSetting.JWTSecretKey)
	if err != nil {
		// 这里不应当出错
		log.Fatal(err)
	}
	return tks
}

func ParseToken(tokenString string, secretKey string) (any, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.ParseWithClaims(tokenString, &RoleClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	return token.Claims, nil
}
