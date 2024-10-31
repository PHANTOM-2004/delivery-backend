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
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func GetAuth(c *gin.Context) {
	// 通过refresh_token, 获得access_token
	//
	account := c.GetString("jwt_account")

	// 提供access_token
	access_token := service.GetAdminAccessToken(account)
	service.SetAccessToken(c, access_token)

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
	// NOTE: add redis blacklist refresh_token
	access_token, err := c.Cookie("access_token")
	if errors.Is(err, http.ErrNoCookie) {
		log.Debug("logout when no access_token provided")
	} else {
		service.DisableAdminToken(access_token, setting.AppSetting.AdminAKAge)
	}

	refresh_token, err := c.Cookie("refresh_token")
	if errors.Is(err, http.ErrNoCookie) {
		log.Debug("logout when no refresh_token provided")
	} else {
		service.DisableAdminToken(refresh_token, setting.AppSetting.AdminRKAge)
	}

	service.DeleteTokens(c)
	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
}

func AdminLogin(c *gin.Context) {
	account := c.PostForm("account")
	password := c.PostForm("password")

	if v := service.AccountValidate(account, password, c); !v {
		return
	}

	// return refresh_token, access_token
	refresh_token := service.GetAdminRefreshToken(account)
	access_token := service.GetAdminAccessToken(account)

	service.SetAccessToken(c, access_token)
	service.SetRefreshToken(c, refresh_token)

	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
}

func AdminChangePassword(c *gin.Context) {
	account := c.GetString("jwt_account")

	// 新密码, 首先进行校验
	origin_new_pwd := c.PostForm("password")
	if v := service.AccountValidate(account, origin_new_pwd, c); !v {
		return
	}

	// Encrypt
	new_password := utils.Encrypt(origin_new_pwd, setting.AppSetting.Salt)
	data := map[string]any{
		"password": new_password,
	}

	err := models.EditAdmin(account, data)
	if err != nil {
		// 在这里edit， 应当保证成功；因为数据库是存在的
		app.Response(c, http.StatusInternalServerError, ecode.ERROR, nil)
		log.Warn("Password Update Failure[internal]")
		return
	}

	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
}
