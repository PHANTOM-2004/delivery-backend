package merchant_service

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
	Account      string
	Password     string
	MerchantName string
	PhoneNumber  string
}

func init() {
	// NOTE: 在该init中，初始化该模块的数据验证
	{
		// register login validation
		// 约定账号长度最大值为30, 最小值为10
		// 约定密码最大长度为32, 最小长度为12

		// 修改密码validation
		password_rules := map[string]string{
			"Password": "min=8,max=30",
		}

		app.RegisterValidation(Password{}, password_rules)

		// 登录validation
		login_rules := password_rules
		login_rules["Account"] = "min=6,max=30"

		app.RegisterValidation(Login{}, login_rules)

		// 注册账号validation
		signup_rules := login_rules
		signup_rules["MerchantName"] = "min=2,max=20"
		// example: +8613912345678
		signup_rules["PhoneNumber"] = "required,e164"
		app.RegisterValidation(SignUp{}, signup_rules)
	}
}

func AccountValidate(account string, password string, c *gin.Context) bool {
	data := Login{
		Account:  account,
		Password: password,
	}

	err := app.ValidateStruct(data)
	if err != nil {
		// 通常来说前端不应当传递非法参数，对于非法参数的传递
		// 通常是其他人所进行的
		log.Warn("Login: invalid params")
		log.Debug(err, data)
		app.Response(c, http.StatusOK, ecode.INVALID_PARAMS, nil)
		return false
	}

	m, err := models.GetMerchant(data.Account)
	if err != nil {
		// 其他未知错误
		log.Warn(err)
		app.Response(c, http.StatusInternalServerError, ecode.ERROR, nil)
		return false
	} else if m == nil {
		// 商家不存在
		log.Debug(err, data)
		app.Response(c, http.StatusOK, ecode.ERROR_MERCHANT_NON_EXIST, nil)
		return false
	}

	en_pwd := utils.Encrypt(data.Password, setting.AppSetting.Salt)
	if en_pwd != m.Password {
		// 用户输错密码
		log.Debug("incorrect password")
		app.Response(c, http.StatusOK, ecode.ERROR_MERCHANT_INCORRECT_PWD, nil)
		return false
	}
	return true
}

func SignUpValidate(c *gin.Context) bool {
	method := c.Request.Method
	if method != "POST" {
		log.Fatal("invalid usage")
		return false
	}

	data := SignUp{
		Account:      c.PostForm("account"),
		Password:     c.PostForm("password"),
		MerchantName: c.PostForm("merchant_name"),
		PhoneNumber:  c.PostForm("phone_number"),
	}

	err := app.ValidateStruct(data)
	if err != nil {
		app.Response(c, http.StatusOK, ecode.INVALID_PARAMS, nil)
		log.Warn(err)
		return false
	}
	return true
}

func PasswordValidate(password string, c *gin.Context) bool {
	data := Password{password}

	err := app.ValidateStruct(data)
	if err != nil {
		app.Response(c, http.StatusOK, ecode.INVALID_PARAMS, nil)
		log.Warn(err)
		return false
	}
	return true
}
