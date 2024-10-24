package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/internal/setting"
	"delivery-backend/models"
	"delivery-backend/pkg/utils"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func AdminChangePassword(c *gin.Context) {
	session := sessions.Default(c)
	account := session.Get("account")
	role := session.Get("role")
	if account == nil || role != "admin" {
		app.Response(c, http.StatusUnauthorized, ecode.ERROR_ADMIN_NOT_LOGIN, nil)
		return
	}

	// Encrypt
	new_password := utils.Encrypt(c.Query("new_password"), setting.AppSetting.Salt)
	data := map[string]any{
		"password": new_password,
	}

	err := models.EditAdmin(account.(string), data)
	if err != nil {
		// 在这里edit， 应当保证成功；因为数据库是存在的
		app.Response(c, http.StatusInternalServerError, ecode.ERROR, nil)
		log.Fatal("Password Update Failure[internal]")
		return
	}

	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
}
