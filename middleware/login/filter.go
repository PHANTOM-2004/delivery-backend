package filter

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LoginFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		account := session.Get("account")

		if account != nil {
			app.Response(c, http.StatusUnauthorized, ecode.ERROR_ADMIN_NOT_LOGIN, nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
