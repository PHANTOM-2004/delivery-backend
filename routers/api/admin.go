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

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Login struct {
	Account  string
	Password string
}

func init() {
	// NOTE: 在该init中，初始化该模块的数据验证
	{
		// register login validation
		// 约定账号长度最大值为30, 最小值为10
		// 约定密码最大长度为32, 最小长度为12
		login_rules := map[string]string{
			"Account":  "min=10,max=30",
			"Password": "min=12,max=50",
		}
		app.RegisterValidation(Login{}, login_rules)
	}
}

func ValidateAccount(c *gin.Context) {
	data := Login{
		Account: c.Query("account"),
		// 密码经过加密
		Password: utils.Encrypt(c.Query("password"), setting.AppSetting.Salt),
	}

	err := app.ValidateStruct(data)
	if err != nil {
		// 通常来说前端不应当传递非法参数，对于非法参数的传递
		// 通常是其他人所进行的
		log.Warn("Login: invalid params")
		log.Debug(err, data)
		app.Response(c, http.StatusOK, ecode.INVALID_PARAMS, nil)
		return
	}

	a, err := models.GetAdmin(data.Account)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 对于不存在的账户登陆，这时可能的，因为
		// 你无法预料到用户会干什么
		log.Debug(err, data)
		app.Response(c, http.StatusOK, ecode.ERROR_ADMIN_NON_EXIST, nil)
		return
	} else if err != nil {
		// 其他未知错误
		log.Warn(err)
		app.Response(c, http.StatusInternalServerError, ecode.ERROR, nil)
		return
	}

	if data.Password != a.Password {
		// 用户输错密码
		log.Debug("incorrect password")
		app.Response(c, http.StatusOK, ecode.ERROR_ADMIN_INCORRECT_PWD, nil)
		return
	}

	// TODO: return session id, and access token
	access_token := service.GetAdminAccessToken(data.Account)

	// NOTE:返回 session_id, access_token
	res := map[string]any{
		"access_token": access_token,
		"session_id":   1,
	}

	app.Response(c, http.StatusOK, ecode.SUCCESS, res)
}
