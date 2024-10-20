package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/models"
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

	account := c.Query("account")
	password := c.Query("password")

	data := map[string]any{
		"password": password,
	}

	err := models.EditAdmin(account, data)
	if err != nil {
		// 在这里edit， 应当保证成功；因为数据库是存在的
		log.Fatal("Password Update Failure[internal]")
		return
	}

	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
}
