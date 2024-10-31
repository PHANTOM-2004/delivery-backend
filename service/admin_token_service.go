package service

import (
	"delivery-backend/internal/ecode"
	"delivery-backend/internal/gredis"
	"delivery-backend/internal/setting"
	"delivery-backend/pkg/utils"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
)

type AdminClaims struct {
	Account string `json:"account"`
	jwt.RegisteredClaims
}

func GetAdminAccessToken(account string) string {
	expires := time.Now().Add(time.Duration(setting.AppSetting.AdminAKAge) * time.Minute)

	claims := AdminClaims{
		account,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
			Issuer:    "admin",
			Subject:   "access_token",
		},
	}

	tks, err := utils.GenerateToken(claims, setting.AppSetting.JWTSecretKey)
	if err != nil {
		// 这里不应当出错
		log.Fatal(err)
	}
	return tks
}

func GetAdminRefreshToken(account string) string {
	expires := time.Now().Add(time.Duration(setting.AppSetting.AdminRKAge) * time.Minute)

	claims := AdminClaims{
		account,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
			Issuer:    "admin",
			Subject:   "refresh_token",
		},
	}

	tks, err := utils.GenerateToken(claims, setting.AppSetting.JWTSecretKey)
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
	token, err := jwt.ParseWithClaims(tokenString, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
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

func AuthAdminAccessToken(access_token string) (string, ecode.Ecode) {
	res_claims, err := ParseToken(access_token, setting.AppSetting.JWTSecretKey)
	if errors.Is(err, jwt.ErrTokenExpired) {
		log.Debug("fail: access_token expired")
		return "", ecode.ERROR_AUTH_ACCESS_TOKEN_EXPIRED
	}

	log.Debug("pass: access_token not expired")

	if err != nil {
		log.Debug(err)
		return "", ecode.ERROR_AUTH_CHECK_ACCESS_TOKEN
	}

	claims, _ := res_claims.(*AdminClaims)

	account := claims.Account
	issuer := claims.Issuer
	sub := claims.Subject

	if issuer != "admin" || sub != "access_token" {
		return "", ecode.ERROR_AUTH_CHECK_ACCESS_TOKEN
	}

	return account, ecode.SUCCESS
}

func AuthAdminRefreshToken(refresh_token string) (string, ecode.Ecode) {
	res_claims, err := ParseToken(refresh_token, setting.AppSetting.JWTSecretKey)

	if errors.Is(err, jwt.ErrTokenExpired) {
		log.Debug("refresh_token expired")
		return "", ecode.ERROR_AUTH_ACCESS_TOKEN_EXPIRED
	}

	log.Debug("pass: refresh_token not expired")

	if err != nil {
		log.Debug(err)
		return "", ecode.ERROR_AUTH_CHECK_REFRESH_TOKEN
	}

	claims, ok := res_claims.(*AdminClaims)
	if !ok {
	}

	account := claims.Account
	issuer := claims.Issuer
	sub := claims.Subject

	if issuer != "admin" || sub != "refresh_token" {
		return "", ecode.ERROR_AUTH_CHECK_REFRESH_TOKEN
	}

	return account, ecode.SUCCESS
}

func DisableAdminToken(admin_token string, expire_minute int) error {
	key := "ADMIN_TK_" + admin_token
	return gredis.Set(key, "", time.Duration(expire_minute)*time.Minute)
}

// 判断token是否在redis黑名单中
// 当token有效返回true，无效返回false
func ValidateAdminToken(admin_token string) bool {
	key := "ADMIN_TK_" + admin_token
	return !gredis.Exists(key)
}

func DeleteTokens(c *gin.Context) {
	// TODO: add redis blacklist refresh_token
	// also add middle for blacklist refresh_token check

	c.SetCookie(
		"access_token",
		"",
		-1,
		"",
		"",
		true,
		true)

	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"",
		"",
		true,
		true)
}

func SetRefreshToken(c *gin.Context, refresh_token string) {
	c.SetCookie(
		"refresh_token",
		refresh_token,
		(setting.AppSetting.AdminRKAge+2)*60,
		"",
		"",
		true,
		true)
}

func SetAccessToken(c *gin.Context, access_token string) {
	c.SetCookie(
		"access_token",
		access_token,
		(setting.AppSetting.AdminAKAge+5)*60,
		"",
		"",
		true,
		true)
}
