package service

import (
	"delivery-backend/internal/ecode"
	"delivery-backend/internal/setting"
	"delivery-backend/models"
	"delivery-backend/pkg/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
)

func GetAccessToken(account string, password string) string {
	claims := jwt.MapClaims{
		"account":    account,
		"password":   password,
		"expires_at": time.Now().Add(3 * time.Hour).Unix(),
	}
	tks, err := utils.GenerateToken(claims, setting.AppSetting.JWTSecretKey)
	if err != nil {
		// 这里不应当出错
		log.Fatal(err)
	}
	return tks
}

func AuthAccessToken(access_token string) ecode.Ecode {
	claims, err := utils.ParseToken(access_token, setting.AppSetting.JWTSecretKey)
	if err != nil {
		return ecode.ERROR_AUTH_CHECK_TOKEN_FAIL
	}
	account := claims["account"].(string)
	password := claims["password"].(string)
	a, err := models.GetAdmin(account)
	if err != nil || a.Password != password {
		return ecode.ERROR_AUTH_CHECK_TOKEN_FAIL
	}

	expires_at := claims["expires_at"].(int64)
	nowTime := time.Now().Unix()
	if nowTime > expires_at {
		return ecode.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
	}

	return ecode.SUCCESS
}
