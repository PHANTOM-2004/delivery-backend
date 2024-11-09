package filter

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// UNUSED
func LoginFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		account := session.Get("account")
		role := session.Get("role")

		if account == nil {
			app.Response(c, http.StatusUnauthorized,
				ecode.ERROR_ADMIN_NOT_LOGIN, nil)
			c.Abort()
			return
		}

		if role != "admin" {
			app.Response(c, http.StatusUnauthorized,
				ecode.ERROR_ADMIN_ROLE, nil)
			c.Abort()
			return
		}

		// set session account
		c.Set("session_account", account)
		// pass
		c.Next()
	}
}

// UNUSED
func DoubleValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		session_account, exist := c.Get("session_account")
		if !exist {
			log.Fatal("account not provided by session")
			app.Response(c, http.StatusInternalServerError, ecode.ERROR, nil)
			c.Abort()
			return
		}

		jwt_account, exist := c.Get("jwt_account")
		if !exist {
			log.Fatal("account not provided by jwt")
			app.Response(c, http.StatusInternalServerError, ecode.ERROR, nil)
			c.Abort()
			return
		}

		if jwt_account != session_account {
			app.Response(c, http.StatusUnauthorized, ecode.ERROR, nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
