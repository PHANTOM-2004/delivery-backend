package service

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/internal/setting"
	"delivery-backend/models"
	"delivery-backend/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
)

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
		login_rules := map[string]string{
			"Account":  "min=10,max=30",
			"Password": "min=15,max=30",
		}
		app.RegisterValidation(Login{}, login_rules)

		// 注册账号validation
		signup_rules := login_rules
		signup_rules["AdminName"] = "min=2,max=20"
		app.RegisterValidation(SignUp{}, signup_rules)
	}
}

func SignUpValidate(c *gin.Context) bool {
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
	err := app.ValidateStruct(data)
	if err != nil {
		app.Response(c, http.StatusOK, ecode.INVALID_PARAMS, nil)
		log.Warn(err)
		return false
	}
	return true
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

	a, err := models.GetAdmin(data.Account)
	if a == nil {
		// 对于不存在的账户登陆，这时可能的，因为
		// 你无法预料到用户会干什么
		log.Debug(err, data)
		app.Response(c, http.StatusOK, ecode.ERROR_ADMIN_NON_EXIST, nil)
		return false
	} else if err != nil {
		// 其他未知错误
		log.Warn(err)
		app.Response(c, http.StatusInternalServerError, ecode.ERROR, nil)
		return false
	}

	en_pwd := utils.Encrypt(data.Password, setting.AppSetting.Salt)
	if en_pwd != a.Password {
		// 用户输错密码
		log.Debug("incorrect password")
		app.Response(c, http.StatusOK, ecode.ERROR_ADMIN_INCORRECT_PWD, nil)
		return false
	}
	return true
}

func GetAdminAccessToken(account string) string {
	claims := jwt.MapClaims{
		"issuer":     "admin",
		"type":       "access",
		"account":    account,
		"expires_at": time.Now().Add(10 * time.Minute).Unix(),
	}
	tks, err := utils.GenerateToken(claims, setting.AppSetting.JWTSecretKey)
	if err != nil {
		// 这里不应当出错
		log.Fatal(err)
	}
	return tks
}

func GetAdminRefreshToken(account string) string {
	claims := jwt.MapClaims{
		"issuer":     "admin",
		"type":       "refresh",
		"account":    account,
		"expires_at": time.Now().Add(25 * time.Minute).Unix(),
	}
	tks, err := utils.GenerateToken(claims, setting.AppSetting.JWTSecretKey)
	if err != nil {
		// 这里不应当出错
		log.Fatal(err)
	}
	return tks
}

func AuthAdminAccessToken(access_token string) (string, ecode.Ecode) {
	claims, err := utils.ParseToken(access_token, setting.AppSetting.JWTSecretKey)
	if err != nil {
		return "", ecode.ERROR_AUTH_CHECK_ACCESS_TOKEN
	}
	account := claims["account"].(string)
	issuer := claims["issuer"].(string)
	t := claims["type"].(string)
	exist, err := models.ExistAdmin(account)
	if err != nil || !exist || issuer != "admin" || t != "access" {
		return "", ecode.ERROR_AUTH_CHECK_ACCESS_TOKEN
	}

	expires_at := claims["expires_at"].(int64)
	nowTime := time.Now().Unix()
	if nowTime > expires_at {
		return "", ecode.ERROR_AUTH_ACCESS_TOKEN_EXPIRED
	}

	return account, ecode.SUCCESS
}

func AuthAdminRefreshToken(refresh_token string) (string, ecode.Ecode) {
	claims, err := utils.ParseToken(refresh_token, setting.AppSetting.JWTSecretKey)
	if err != nil {
		return "", ecode.ERROR_AUTH_CHECK_REFRESH_TOKEN
	}
	account := claims["account"].(string)
	issuer := claims["issuer"].(string)
	t := claims["type"].(string)
	exist, err := models.ExistAdmin(account)
	if err != nil || !exist || issuer != "admin" || t != "refresh" {
		return "", ecode.ERROR_AUTH_CHECK_REFRESH_TOKEN
	}

	expires_at := claims["expires_at"].(int64)
	nowTime := time.Now().Unix()
	if nowTime > expires_at {
		return "", ecode.ERROR_AUTH_REFRESH_TOKEN_EXPIRED
	}

	return account, ecode.SUCCESS
}

func DeleteTokens(c *gin.Context) {
	c.SetCookie(
		"access_token",
		"",
		-1,
		"",
		"",
		true,
		true)

	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"",
		"",
		true,
		true)
}

func SetRefreshToken(c *gin.Context, refresh_token string) {
	c.SetCookie(
		"refresh_token",
		refresh_token,
		(setting.AppSetting.AdminRKAge+2)*60,
		"",
		"",
		true,
		true)
}

func SetAccessToken(c *gin.Context, access_token string) {
	c.SetCookie(
		"access_token",
		access_token,
		(setting.AppSetting.AdminAKAge+5)*60,
		"",
		"",
		true,
		true)
}
