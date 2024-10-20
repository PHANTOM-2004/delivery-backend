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
	"delivery-backend/service"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// NOTE:该接口只允许运维调用，需要验证创建管理员的唯一token.
func AdminCreate(c *gin.Context) {

}

// NOTE:我们要求前端使用cookie方式存储access_token，
// 并且设置secure，https only
func GetAuth(c *gin.Context) {
	if v := service.AccountValidate(c); !v {
		return
	}

	// 提供access_token
	account := c.Query("account")
	access_token := service.GetAdminAccessToken(account)
	res := map[string]any{
		"access_token": access_token,
	}
	app.Response(c, http.StatusOK, ecode.SUCCESS, res)
}

func AdminLogin(c *gin.Context) {
	if v := service.AccountValidate(c); !v {
		return
	}

	account := c.Query("account")

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
