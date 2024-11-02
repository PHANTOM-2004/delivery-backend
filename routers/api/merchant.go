package api

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/internal/setting"
	"delivery-backend/models"
	"delivery-backend/pkg/utils"
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

func MerchantLogin(c *gin.Context) {
	account := c.PostForm("account")
	password := c.PostForm("password")

	if v := merchant_service.AccountValidate(account, password, c); !v {
		return
	}

	// return refresh_token, access_token
	refresh_token := merchant_service.GetRefreshToken(account)
	access_token := merchant_service.GetAccessToken(account)

	merchant_service.SetAccessToken(c, access_token)
	merchant_service.SetRefreshToken(c, refresh_token)

	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
}

func MerchantChangePassword(c *gin.Context) {
	account := c.GetString("jwt_account")

	new_pwd := c.PostForm("password")
	// 新密码, 首先进行校验
	if v := merchant_service.PasswordValidate(new_pwd, c); !v {
		return
	}

	// Encrypt
	new_password := utils.Encrypt(new_pwd, setting.AppSetting.Salt)
	data := map[string]any{
		"password": new_password,
	}

	id, err := models.GetMerchantID(account)
	if err == nil {
		err = models.EditMerchant(id, data)
	}

	if err != nil {
		// 在这里edit， 应当保证成功；因为数据库是存在的
		app.Response(c, http.StatusInternalServerError, ecode.ERROR, nil)
		log.Warn("Password Update Failure[internal]")
		return
	}

	// 应当删除tokens
	merchant_service.DisbleTokens(c)

	// 返回响应
	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
	log.Debug("password updated")
}


