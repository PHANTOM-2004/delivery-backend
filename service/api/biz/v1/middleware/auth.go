package middleware

import (
	"delivery-backend/common/app"
	"delivery-backend/common/ecode"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		warp := app.RespWarp{Context: c}
		session := sessions.Default(c)
		admin_id := session.Get("admin_id")
		if admin_id == nil {
			warp.Resp(http.StatusUnauthorized,
				ecode.ERROR_ADMIN_NOT_LOGIN, nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

func MerchAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		warp := app.RespWarp{Context: c}
		session := sessions.Default(c)
		merch_id := session.Get("merchant_id")
		if merch_id == nil {
			warp.Resp(http.StatusUnauthorized,
				ecode.ERROR_MERCHANT_NOT_LOGIN, nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
