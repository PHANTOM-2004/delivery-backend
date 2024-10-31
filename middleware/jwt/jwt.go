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

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		access_token, err := c.Cookie("access_token")
		if errors.Is(err, http.ErrNoCookie) {
			// TODO:
			app.Response(c, http.StatusOK, ecode.ERROR_AUTH_NO_ACCESS_TOKEN, nil)
			c.Abort()
			return
		}

		account, code := service.AuthAdminAccessToken(access_token)
		if code != ecode.SUCCESS {
			app.Response(c, http.StatusOK, code, nil)
			c.Abort()
			return
		}

		c.Set("jwt_account", account)
		c.Next()

		log.Debug("pass jwt access_token")
	}
}
