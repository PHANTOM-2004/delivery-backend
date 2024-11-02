package api

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/service/merchant_service"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func MerchantGetAuth(c *gin.Context) {
	// 通过refresh_token, 获得access_token
	// 通过refresh_token, 获得access_token
	//
	account := c.GetString("jwt_account")

	// 提供access_token
	access_token := merchant_service.GetAccessToken(account)
	merchant_service.SetAccessToken(c, access_token)

	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)

	log.Debug("pass: response access_token")
	// account := c.GetString("jwt_account")
}

func MerchantLoginStatus(c *gin.Context) {
	account := c.GetString("jwt_account")
	res := map[string]string{
		"account": account,
	}
	app.Response(c, http.StatusOK, ecode.SUCCESS, res)
}

func MerchantLogout(c *gin.Context) {
	merchant_service.DisbleTokens(c)
	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
}
