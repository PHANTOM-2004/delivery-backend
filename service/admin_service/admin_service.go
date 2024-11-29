package admin_service

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/internal/setting"
	"delivery-backend/models"
	"delivery-backend/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Password struct {
	Password string
}

type Login struct {
	Account  string
	Password string
}

type SignUp struct {
	Account   string
	Password  string
	AdminName string
}

func init() {
	// NOTE: 在该init中，初始化该模块的数据验证
	{
		// register login validation
		// 约定账号长度最大值为30, 最小值为10
		// 约定密码最大长度为32, 最小长度为12

		// 修改密码validation
		password_rules := map[string]string{
			"Password": "min=15,max=30",
		}

		app.RegisterMapValidation(Password{}, password_rules)

		// 登录validation
		login_rules := password_rules
		login_rules["Account"] = "min=10,max=30"

		app.RegisterMapValidation(Login{}, login_rules)

		// 注册账号validation
		signup_rules := login_rules
		signup_rules["AdminName"] = "min=2,max=20"
		app.RegisterMapValidation(SignUp{}, signup_rules)
	}
}

func SignUpRequestValidate(c *gin.Context) bool {
	method := c.Request.Method
	if method != "POST" {
		log.Fatal("invalid usage")
		return false
	}

	data := SignUp{
		Account:   c.Query("account"),
		Password:  c.Query("password"),
		AdminName: c.Query("admin_name"),
	}
	err := app.ValidateStruct(&data)
	if err != nil {
		app.ResponseInvalidParams(c)
		log.Warn(err)
		return false
	}
	return true
}

func PasswordRequestValidate(password string, c *gin.Context) bool {
	data := Password{password}

	err := app.ValidateStruct(&data)
	if err != nil {
		app.ResponseInvalidParams(c)
		log.Warn(err)
		return false
	}
	return true
}

// 账号登陆验证
func AdminLoginValidate(account string, password string, c *gin.Context) (uint, bool) {
	data := Login{
		Account:  account,
		Password: password,
	}

	err := app.ValidateStruct(&data)
	if err != nil {
		// 通常来说前端不应当传递非法参数，对于非法参数的传递
		// 通常是其他人所进行的
		log.Warn(err, data)
		app.ResponseInvalidParams(c)
		return 0, false
	}

	a, err := models.GetAdmin(data.Account)
	if err != nil {
		// 其他未知错误
		app.ResponseInternalError(c, err)
		return 0, false
	} else if a.ID == 0 {
		// 对于不存在的账户登陆，这时可能的，因为
		// 你无法预料到用户会干什么
		log.Debug(err, data)
		app.Response(c, http.StatusOK, ecode.ERROR_ADMIN_NOT_FOUND, nil)
		return 0, false

	}

	en_pwd := utils.Encrypt(data.Password, setting.AppSetting.Salt)
	if en_pwd != a.Password {
		// 用户输错密码
		log.Debug("incorrect password")
		app.Response(c, http.StatusOK, ecode.ERROR_ADMIN_INCORRECT_PWD, nil)
		return 0, false
	}
	return a.ID, true
}
