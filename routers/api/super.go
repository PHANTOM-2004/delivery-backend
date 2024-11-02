package api

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/internal/setting"
	"delivery-backend/models"
	"delivery-backend/pkg/utils"
	"delivery-backend/service/admin_service"
	"net/http"

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
	if v := admin_service.SignUpValidate(c); !v {
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
