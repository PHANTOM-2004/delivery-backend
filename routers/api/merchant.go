package api

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/internal/setting"
	"delivery-backend/models"
	"delivery-backend/pkg/utils"
	handler "delivery-backend/service"
	"delivery-backend/service/merchant_service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func DeleteMerchant(c *gin.Context) {
	account := c.PostForm("account")
	err, rows := models.DeleteMerchant(account)
	if err != nil {
		res := map[string]string{
			"error": err.Error(),
		}
		app.Response(c, http.StatusOK, ecode.ERROR, res)
		return
	}

	if rows <= 0 {
		res := map[string]string{
			"warn": "delete nothing",
		}
		app.Response(c, http.StatusOK, ecode.SUCCESS, res)
		return
	}

	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
}

func CreateMerchant(c *gin.Context) {
	if v := merchant_service.SignUpRequestValidate(c); !v {
		return
	}

	// create account
	account := c.PostForm("account")
	encrypted_password := utils.Encrypt(c.PostForm("password"), setting.AppSetting.Salt)
	merchant_name := c.PostForm("merchant_name")
	phone_numer := c.PostForm("phone_number")
	application_id, err := strconv.Atoi(c.PostForm("merchant_application_id"))
	if err != nil {
		app.ResponseInvalidParams(c)
		return
	}

	data := models.Merchant{
		MerchantName:          merchant_name,
		Account:               account,
		Password:              encrypted_password,
		PhoneNumber:           phone_numer,
		MerchantApplicationID: uint(application_id),
	}

	exist, err := models.ExistMerchant(account)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	if exist {
		app.Response(c, http.StatusOK, ecode.ERROR_MERCHANT_ACCOUNT_EXIST, nil)
		return
	}

	err = models.CreateMerchant(&data)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
}

func MerchantLoginStatus(c *gin.Context) {
	account := handler.NewMerchInfoHanlder(c).GetAccount()
	res := map[string]string{
		"account": account,
	}
	app.Response(c, http.StatusOK, ecode.SUCCESS, res)
}

func MerchantLogout(c *gin.Context) {
	h := handler.NewMerchInfoHanlder(c)
	err := h.Delete()
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
}

// 认证成功时会在该函数中设置
// c.Set("merchant_id", id)
func MerchantLogin(c *gin.Context) {
	account := c.PostForm("account")
	password := c.PostForm("password")

	id, v := merchant_service.AccountAuth(account, password, c)
	if !v {
		return
	}

	h := handler.NewMerchInfoHanlder(c)
	h.SetID(id)
	h.SetAccount(account)
	err := h.Save()
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
}

func MerchantChangePassword(c *gin.Context) {
	h := handler.NewMerchInfoHanlder(c)
	id := h.GetID()
	new_pwd := c.PostForm("password")
	// 新密码, 首先进行校验
	if v := merchant_service.PasswordRequestValidate(new_pwd, c); !v {
		return
	}

	// Encrypt
	new_password := utils.Encrypt(new_pwd, setting.AppSetting.Salt)
	data := map[string]any{
		"password": new_password,
	}

	err := models.UpdateMerchant(id, data)
	if err != nil {
		// 在这里edit， 应当保证成功；因为数据库是存在的
		app.Response(c, http.StatusInternalServerError, ecode.ERROR, nil)
		log.Warn("Password Update Failure[internal]")
		return
	}

	// 应当删除tokens
	err = h.Delete()
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	// 返回响应
	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
	log.Debug("password updated")
}
