package auth

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	handler "delivery-backend/service"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func MerchantAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := handler.NewMerchInfoHanlder(c)
		account := h.GetAccount()
		id := h.GetID()
		log.Trace(account, id)

		if id == 0 || account == "" {
			app.Response(c, http.StatusUnauthorized, ecode.ERROR_ADMIN_NOT_LOGIN, nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
