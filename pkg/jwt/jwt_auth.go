package jwt_token

import (
	"delivery-backend/internal/ecode"
	"delivery-backend/internal/gredis"
	"delivery-backend/internal/setting"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
)

func AuthAccessToken(issuer string, access_token string) (string, ecode.Ecode) {
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

	claims, _ := res_claims.(*RoleClaims)

	account := claims.Account
	iss := claims.Issuer
	sub := claims.Subject

	if iss != issuer || sub != "access_token" {
		return "", ecode.ERROR_AUTH_CHECK_ACCESS_TOKEN
	}

	return account, ecode.SUCCESS
}

func AuthRefreshToken(issuer string, refresh_token string) (string, ecode.Ecode) {
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

	claims, ok := res_claims.(*RoleClaims)
	if !ok {
	}

	account := claims.Account
	iss := claims.Issuer
	sub := claims.Subject

	if iss != issuer || sub != "refresh_token" {
		return "", ecode.ERROR_AUTH_CHECK_REFRESH_TOKEN
	}

	return account, ecode.SUCCESS
}

func deleteTokens(c *gin.Context) {
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

	log.Debug("tokens clears")
}

func disableToken(prefix string, token string, expire_minute int) error {
	key := prefix + token
	return gredis.Set(key, "", time.Duration(expire_minute)*time.Minute)
}

// 判断token是否在redis黑名单中
// 当token有效返回true，无效返回false
func ValidateToken(prefix string, token string) bool {
	key := prefix + token
	return !gredis.Exists(key)
}

func DisableTokens(prefix string, c *gin.Context) {
	// NOTE: add redis blacklist refresh_token
	access_token, err := c.Cookie("access_token")
	if errors.Is(err, http.ErrNoCookie) {
		log.Debug("logout when no access_token provided")
	} else {
		disableToken(prefix, access_token, setting.AppSetting.AdminAKAge)
	}

	refresh_token, err := c.Cookie("refresh_token")
	if errors.Is(err, http.ErrNoCookie) {
		log.Debug("logout when no refresh_token provided")
	} else {
		disableToken(prefix, refresh_token, setting.AppSetting.AdminRKAge)
	}

	deleteTokens(c)
}

func SetRefreshToken(c *gin.Context, refresh_token string, expire_minute int) {
	c.SetCookie(
		"refresh_token",
		refresh_token,
		(expire_minute)*60,
		"",
		"",
		true,
		true)
}

func SetAccessToken(c *gin.Context, access_token string, expire_minute int) {
	c.SetCookie(
		"access_token",
		access_token,
		(expire_minute)*60,
		"",
		"",
		true,
		true)
}
