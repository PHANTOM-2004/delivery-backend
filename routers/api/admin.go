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
	"delivery-backend/service"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetAuth(c *gin.Context) {
	account := c.Query("account")
	password := c.Query("password")
	if v := service.AccountValidate(account, password, c); !v {
		return
	}

	// 提供access_token
	access_token := service.GetAdminAccessToken(account)
	res := map[string]any{
		"access_token": access_token,
	}
	app.Response(c, http.StatusOK, ecode.SUCCESS, res)
}

func AdminLoginStatus(c *gin.Context) {
	account, _ := c.Get("session_account")
	res := map[string]string{
		"account": account.(string),
	}
	app.Response(c, http.StatusOK, ecode.SUCCESS, res)
}

func AdminLogout(c *gin.Context) {
	session := sessions.Default(c)
	// this will mark the session as "written" only if there's
	// at least one key to delete
	session.Clear() // account to delete
	session.Options(sessions.Options{MaxAge: -1})
	session.Save()
	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
}

func AdminLogin(c *gin.Context) {
	session := sessions.Default(c)
	session_account := session.Get("account")
	session_role := session.Get("role")
	if session_account != nil && session_role == "admin" {
		// 如果已经登陆，那么利用session的期限直接认证即可
		app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
		return
	}

	account := c.PostForm("account")
	password := c.PostForm("password")

	if v := service.AccountValidate(account, password, c); !v {
		return
	}

	// set account and role
	session.Set("account", account)
	session.Set("role", "admin")
	session.Save()

	// NOTE:
	// 不必考虑session id的问题，gin-session作为中间件自动管理session,
	// 但是需要注意，这里的自动管理是基于cookie的，也就是说
	// 对于之后的小程序业务，就不能使用gin-session了

	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
}

func AdminChangePassword(c *gin.Context) {
	account, _ := c.Get("session_account")

	// 新密码, 首先进行校验
	origin_new_pwd := c.PostForm("password")
	if v := service.AccountValidate(account.(string), origin_new_pwd, c); !v {
		return
	}

	// Encrypt
	new_password := utils.Encrypt(origin_new_pwd, setting.AppSetting.Salt)
	data := map[string]any{
		"password": new_password,
	}

	err := models.EditAdmin(account.(string), data)
	if err != nil {
		// 在这里edit， 应当保证成功；因为数据库是存在的
		app.Response(c, http.StatusInternalServerError, ecode.ERROR, nil)
		log.Warn("Password Update Failure[internal]")
		return
	}

	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
}
