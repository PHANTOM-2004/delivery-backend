package wechat

import (
	"delivery-backend/internal/app"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type WXSessionInfo struct {
	openid int
}

// 从中提取出session_id
func WXsession() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO:
		// handle 微信发送的session id 请求
		session_id := c.PostForm("session_id")
		if session_id == "" {
			log.Debug("没有提供session_id")
			app.ResponseInvalidParams(c)
			return
		}

		c.Next()
	}
}
