package api

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	handler "delivery-backend/service"
	"delivery-backend/service/merchant_service"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

	id, v := merchant_service.MerchantLoginValidate(account, password, c)
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
