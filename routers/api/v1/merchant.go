package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/internal/setting"
	"delivery-backend/models"
	"delivery-backend/pkg/utils"
	handler "delivery-backend/service"
	"delivery-backend/service/merchant_service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func MerchantLoginStatus(c *gin.Context) {
	account := handler.NewMerchInfoHanlder(c).GetAccount()
	res := map[string]string{
		"account": account,
	}
	app.Response(c, http.StatusOK, ecode.SUCCESS, res)
}

// 录入店铺开启时间
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

func GetMerchantInfo(c *gin.Context) {
	id := handler.NewMerchInfoHanlder(c).GetID()

	m, err := models.GetMerchantByID(id)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	if m == nil {
		// 不可能发生
		err := fmt.Errorf("merchant id[%v] not exist，but this should not happen", id)
		app.ResponseInternalError(c, err)
		log.Panic(err)
		return
	}

	data := map[string]any{
		"merchant": *m,
	}
	app.Response(c, http.StatusOK, ecode.SUCCESS, data)
}
