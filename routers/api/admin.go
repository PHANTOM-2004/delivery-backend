package api

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/internal/setting"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/scrypt"
)

func Crypt(s string) string {
	// DO NOT use this salt value; generate your own random salt. 8 bytes is
	// a good length.
	salt := []byte(setting.AppSetting.Salt)

	dk, err := scrypt.Key([]byte(s), salt, 1<<15, 8, 1, 32)
	if err != nil {
		log.Fatal(err)
	}
	res := base64.StdEncoding.EncodeToString(dk)
	return res
}

type Login struct {
	Account  string
	Password string
}

func init() {
	// NOTE: 在该init中，初始化该模块的数据验证
	{
		rules := map[string]string{
			"Account":  "min=10,max=30",
			"Password": "min=12,max=32",
		}
		app.RegisterValidation(Login{}, rules)
	}
}

func ValidateAccount(c *gin.Context) {
	data := Login{
		Account:  c.Query("account"),
		Password: c.Query("password"),
	}

	err := app.ValidateStruct(data)
	if err != nil {
		app.Response(c, http.StatusOK, ecode.INVALID_PARAMS, nil)
		log.Warn("Login: invalid params")
		return
	}

	res := make(map[string]any)
	// TODO: return session id, and access token
	//
	app.Response(c, http.StatusOK, ecode.SUCCESS, res)

	return
}
