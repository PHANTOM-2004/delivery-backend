package jwt

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		access_token := c.Query("access_token")
		if access_token == "" {
			app.Response(c, http.StatusOK, ecode.ERROR_AUTH_NO_TOKEN, nil)
			c.Abort()
			return
		}

		account, code := service.AuthAdminAccessToken(access_token)

		if code != ecode.SUCCESS {
			app.Response(c, http.StatusOK, code, nil)
			c.Abort()
			return
		}

		// SUCCESS, set account
		c.Set("account", account)
		// go to next handler
		c.Next()
	}
}
