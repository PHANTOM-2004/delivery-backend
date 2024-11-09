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
	"delivery-backend/middleware/jwt"
	"delivery-backend/models"
	"delivery-backend/pkg/utils"
	"delivery-backend/service/admin_service"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func AdminGetAuth(c *gin.Context) {
	// 通过refresh_token, 获得access_token
	//
	account := c.GetString("jwt_account")
	id := c.GetUint("jwt_id")

	// 提供access_token
	access_token := admin_service.GetAccessToken(uint(id), account)
	admin_service.SetAccessToken(c, access_token)

	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)

	log.Debug("pass: response access_token")
}

func AdminLoginStatus(c *gin.Context) {
	account := c.GetString("jwt_account")
	res := map[string]string{
		"account": account,
	}
	app.Response(c, http.StatusOK, ecode.SUCCESS, res)
}

func AdminLogout(c *gin.Context) {
	admin_service.DisbleTokens(c)
	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
}

func AdminLogin(c *gin.Context) {
	account := c.PostForm("account")
	password := c.PostForm("password")
	id, v := admin_service.AccountValidate(account, password, c)
	if !v {
		return
	}

	// return refresh_token, access_token
	refresh_token := admin_service.GetRefreshToken(id, account)
	access_token := admin_service.GetAccessToken(id, account)

	admin_service.SetAccessToken(c, access_token)
	admin_service.SetRefreshToken(c, refresh_token)
	log.Trace("admin tokens set")

	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
}

func AdminChangePassword(c *gin.Context) {
	id := jwt.NewJwtInfo(c).GetID()

	new_pwd := c.PostForm("password")
	// 新密码, 首先进行校验
	if v := admin_service.PasswordValidate(new_pwd, c); !v {
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
		log.Warn(err)
		app.Response(c, http.StatusInternalServerError, ecode.ERROR, nil)
		return
	}

	// 应当删除tokens
	admin_service.DisbleTokens(c)

	// 返回响应
	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
	log.Debug("password updated")
}
