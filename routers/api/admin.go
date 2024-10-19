package api

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/internal/setting"
	"delivery-backend/models"
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/scrypt"
	"gorm.io/gorm"
)

func Encrypt(s string) string {
	// salt: 8 bytes is
	// a good length.
	salt := []byte(setting.AppSetting.Salt)

	dk, err := scrypt.Key([]byte(s), salt, 1<<15, 8, 1, 32)
	if err != nil {
		log.Panic(err)
	}
	res := base64.StdEncoding.EncodeToString(dk)
	return res
}

type Login struct {
	Account  string `json:"admin_account"`
	Password string `json:"admin_password"`
}

func init() {
	// NOTE: 在该init中，初始化该模块的数据验证
	{
		// register login validation
		// 约定账号长度最大值为30, 最小值为10
		// 约定密码最大长度为32, 最小长度为12
		rules := map[string]string{
			"Account":  "min=10,max=30",
			"Password": "min=12,max=50",
		}
		app.RegisterValidation(Login{}, rules)
	}
}

func ExistAccount(c *gin.Context){
  data := c.Query("account")



}

func ValidateAccount(c *gin.Context) {
	data := Login{
		Account: c.Query("account"),
		// 密码经过加密
		Password: Encrypt(c.Query("password")),
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

	// NOTE:返回 session_id, access_token
	res := make(map[string]any)
	res["access_token"] = "TODO"
	res["session_id"] = "TODO"
	// TODO: return session id, and access token
	//
	app.Response(c, http.StatusOK, ecode.SUCCESS, res)
}
