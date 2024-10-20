package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/internal/setting"
	"delivery-backend/models"
	"delivery-backend/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func AdminChangePassword(c *gin.Context) {
	// 这里借助JWT进行鉴权，在中间件中，对于通过验证的请求
	// 会在gin.Context之中设置 account
	account := c.Query("account")
	if account == "" {
		app.Response(c, http.StatusInternalServerError, ecode.ERROR, nil)
		log.Fatal("AdminChangePassword: account cannot be null string")
    return
	}

	// Encrypt
	new_password := utils.Encrypt(c.Query("new_password"), setting.AppSetting.Salt)
	data := map[string]any{
		"password": new_password,
	}

	err := models.EditAdmin(account, data)
	if err != nil {
		// 在这里edit， 应当保证成功；因为数据库是存在的
		app.Response(c, http.StatusInternalServerError, ecode.ERROR, nil)
		log.Fatal("Password Update Failure[internal]")
		return
	}

	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
}
