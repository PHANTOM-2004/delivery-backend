package merchant_service

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/internal/gredis"
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

// 返回redis黑名单中account对应的key
func getAccountKey(account string) string {
	return "MERCH_ACC_" + account
}

func EnableAccount(merchant_id uint, account string) error {
	key := getAccountKey(account)

	exist := gredis.Exists(key)
	if !exist {
		// NOTE: 这不允许发生，对接的前端心里有点B数
		log.Warnf("Enable account[%s] that is not in blacklist", account)
		return nil
	}

	// 写入redis
	err := gredis.Delete(key)
	if err != nil {
		return err
	}

	// 写入数据库
	err = models.EnableMerchant(merchant_id)
	if err != nil {
		return err
	}

	return nil
}

func DisableAccount(merchant_id uint, account string) error {
	key := getAccountKey(account)
	// 首先设置缓存中账户状态为禁用
	err := gredis.Set(key, "", 0)
	if err != nil {
		return err
	}
	// 然后设置数据库商家的状态为禁用状态
	err = models.DisableMerchant(merchant_id)
	if err != nil {
		return err
	}
	return nil
}

// 只从redis中查询状态
func AccountInBlacklist(account string) bool {
	key := getAccountKey(account)
	exist := gredis.Exists(key)
	return exist
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

// 管理员端为Merchant创建
func CreateMerchantFromApplication(application_id uint) error {
	a, err := models.
		GetMerchantApplication(int(application_id))
	if err != nil {
		log.Warnf("Merchant Application id[%d] not found", application_id)
		return err
	}

	// TODO: 暂定注册规则为随机字符串,后续按照需要更改
	account := "M" + utils.RandString(10)
	password := "P" + utils.RandString(11)
	en_password := utils.Encrypt(password, setting.AppSetting.Salt)

	m := models.Merchant{
		Account:               account,
		Password:              en_password,
		PhoneNumber:           a.PhoneNumber,
		MerchantName:          a.Name,
		MerchantApplicationID: int(application_id),
	}

	err = models.CreateMerchant(&m)
	if err == nil {
		log.Debugf("created merchant:account[%s],password[%s]",
      account, password)
	}
	return err
}
