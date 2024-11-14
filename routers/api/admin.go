// 本文件负责编写管理员账户相关的api, 主要是登入认证与密码修改
//
// NOTE:对于管理员的注册，不应当放在用户界面。实际上应当通过
// 后台插入的方式建立管理员账户。因此对于管理员账户，我们仅仅设计
// 登陆，以及修改密码的功能。同时还需要注意对于密码的修改，
// 要保证必然是强密码,为了保证安全性，在后端需要再次对其进行验证，
// 因为我们不能避免管理员/hacker通过api的方式调用并且传递弱密码

package api

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/internal/setting"
	"delivery-backend/models"
	"delivery-backend/pkg/utils"
	handler "delivery-backend/service"
	"delivery-backend/service/admin_service"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func AdminLoginStatus(c *gin.Context) {
	h := handler.NewAdminInfoHanlder(c)
	account := h.GetAccount()
	res := map[string]string{
		"account": account,
	}
	app.Response(c, http.StatusOK, ecode.SUCCESS, res)
}

func AdminLogout(c *gin.Context) {
	h := handler.NewAdminInfoHanlder(c)
	err := h.Delete()
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
}

func AdminLogin(c *gin.Context) {
	account := c.PostForm("account")
	password := c.PostForm("password")
	id, v := admin_service.AdminLoginValidate(account, password, c)
	if !v {
		return
	}

	h := handler.NewAdminInfoHanlder(c)
	h.SetAccount(account)
	h.SetID(id)
	err := h.Save()
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
}

func AdminChangePassword(c *gin.Context) {
	h := handler.NewAdminInfoHanlder(c)
	id := h.GetID()
	new_pwd := c.PostForm("password")
	// 新密码, 首先进行校验
	if v := admin_service.PasswordRequestValidate(new_pwd, c); !v {
		return
	}

	// Encrypt
	new_password := utils.Encrypt(new_pwd, setting.AppSetting.Salt)
	data := map[string]any{
		"password": new_password,
	}

	err := models.EditAdmin(id, data)
	if err != nil {
		// 在这里edit， 应当保证成功；因为数据库是存在的
		log.Warn("Password Update Failure[internal]")
		app.ResponseInternalError(c, err)
		return
	}

	// 删除当前session
	h.Delete()

	// 返回响应
	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
	log.Debug("password updated")
}
