package jwt

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func JWTAK() gin.HandlerFunc {
	return func(c *gin.Context) {
		access_token, err := c.Cookie("access_token")
		if errors.Is(err, http.ErrNoCookie) {
			app.Response(c, http.StatusOK, ecode.ERROR_AUTH_NO_ACCESS_TOKEN, nil)
			c.Abort()
      log.Debug("fail: no access_token provided")
			return
		}

    log.Debug("pass: received jwt access_token")

		account, code := service.AuthAdminAccessToken(access_token)
		if code != ecode.SUCCESS {
			app.Response(c, http.StatusOK, code, nil)
			c.Abort()
			return
		}

		log.Debug("pass: access_token validation")

		// 检查是否在黑名单中
		valid := service.ValidateAdminToken(access_token)
		if !valid {
			app.Response(c, http.StatusOK, ecode.ERROR_AUTH_ACCESS_TOKEN_EXPIRED, nil)
			c.Abort()
			return
		}

		c.Set("jwt_account", account)

		log.Debug("pass: access_token not in blacklist")

		c.Next()
	}
}
