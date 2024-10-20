package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/internal/setting"
	"delivery-backend/models"
	"delivery-backend/pkg/utils"
	"delivery-backend/service"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func AdminChangePassword(c *gin.Context) {
	// TODO:当前账户实际上从session中获取，对于没有登陆的管理员
	// 显然不让他修改密码

	if v := service.AccountValidate(c); !v {
		return
	}

	// Encrypt
	account := c.Query("account")
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
