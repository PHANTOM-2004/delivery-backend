package jwt

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/service/admin_service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func JWTRK() gin.HandlerFunc {
	return func(c *gin.Context) {
		// cookie中读取refresh_token
		refresh_token, err := c.Cookie("refresh_token")
		if errors.Is(err, http.ErrNoCookie) {
			app.Response(c, http.StatusOK, ecode.ERROR_AUTH_NO_REFRESH_TOKEN, nil)
			return
		}

		log.Debug("pass: refresh_token exists")

		// 校验refresh token
		account, code := admin_service.AuthAdminRefreshToken(refresh_token)
		if code != ecode.SUCCESS {
			app.Response(c, http.StatusOK, code, nil)
			return
		}

		log.Debug("pass: refresh_token validation")

		// 检查refresh token是否在黑名单中
		valid := admin_service.ValidateAdminToken(refresh_token)
		if !valid {
			app.Response(c, http.StatusOK, ecode.ERROR_AUTH_REFRESH_TOKEN_EXPIRED, nil)
			c.Abort()
			log.Debug("fail: refresh_token in blacklist")
			return
		}

		log.Debug("pass: refresh_token not in blacklist")

		c.Set("jwt_account", account)

		c.Next()
	}
}

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

		account, code := admin_service.AuthAdminAccessToken(access_token)
		if code != ecode.SUCCESS {
			app.Response(c, http.StatusOK, code, nil)
			c.Abort()
			return
		}

		log.Debug("pass: access_token validation")

		// 检查是否在黑名单中
		valid := admin_service.ValidateAdminToken(access_token)
		if !valid {
			app.Response(c, http.StatusOK, ecode.ERROR_AUTH_ACCESS_TOKEN_EXPIRED, nil)
			c.Abort()
			log.Debug("fail: access_token in blacklist")
			return
		}

		c.Set("jwt_account", account)

		log.Debug("pass: access_token not in blacklist")

		c.Next()
	}
}
