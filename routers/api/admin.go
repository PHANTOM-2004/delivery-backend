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

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func superTokenCheck(c *gin.Context) bool {
	// TODO: 如果启用管理员插入，token应当每30分钟动态生成一次，所以不应当写成硬编码
	// 会把动态生成的token发送给管理员或者写入安全性高的地方。token生成可以通过调整
	// 配置文件关闭掉。
	// 暂时先采用如下方法使用
	if setting.AppSetting.AdminToken == "" {
		// 如果token为空那么意味着禁用
		app.Response(c, http.StatusNotFound, ecode.ERROR, nil)
		return false
	}
	token := c.Query("super_token")
	if token == "" {
		app.Response(c, http.StatusOK, ecode.ERROR_SUPER_AUTH_NO_TOKEN, nil)
		return false
	}
	if token != setting.AppSetting.AdminToken {
		app.Response(c, http.StatusOK, ecode.ERROR_SUPER_AUTH, nil)
		return false
	}
	return true
}

func AdminDelete(c *gin.Context) {
	if v := superTokenCheck(c); !v {
		return
	}

	account := c.Query("account")
	err, rows := models.DeleteAdmin(account)
	if err != nil {
		res := map[string]string{
			"error": err.Error(),
		}
		app.Response(c, http.StatusOK, ecode.ERROR, res)
		return
	}

	if rows <= 0 {
		res := map[string]string{
			"warn": "delete nothing",
		}
		app.Response(c, http.StatusOK, ecode.SUCCESS, res)
		return
	}

	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
}

// NOTE:该接口只允许运维调用，需要验证创建管理员的唯一token.
func AdminCreate(c *gin.Context) {
	if v := superTokenCheck(c); !v {
		return
	}
	if v := service.SignUpValidate(c); !v {
		return
	}

	// create account
	account := c.Query("account")
	encrypted_password := utils.Encrypt(c.Query("password"), setting.AppSetting.Salt)
	admin_name := c.Query("admin_name")

	data := models.Admin{
		AdminName: admin_name,
		Account:   account,
		Password:  encrypted_password,
	}
	err := models.CreateAdmin(&data)
	if err != nil {
		res := map[string]string{
			"error": err.Error(),
		}
		app.Response(c, http.StatusOK, ecode.ERROR, res)
		return
	}

	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
}

// NOTE:我们要求前端使用cookie方式存储access_token，
// 并且设置secure，https only
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

func AdminLogin(c *gin.Context) {
	account := c.PostForm("account")
	password := c.PostForm("password")

	if v := service.AccountValidate(account, password, c); !v {
		return
	}

	// set in JWT middleware
	account = c.MustGet("account").(string)

	session := sessions.Default(c)
	session.Set("account", account)
	session.Save()

	access_token := service.GetAdminAccessToken(account)
	// NOTE:返回 access_token
	// 不必考虑session id的问题，gin-session作为中间件自动管理session,
	// 但是需要注意，这里的自动管理是基于cookie的，也就是说
	// 对于之后的小程序业务，就不能使用gin-session了
	res := map[string]any{
		"access_token": access_token,
	}

	app.Response(c, http.StatusOK, ecode.SUCCESS, res)
}
