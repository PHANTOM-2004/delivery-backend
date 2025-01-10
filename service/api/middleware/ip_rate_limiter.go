package middleware

import (
	"delivery-backend/common/app"
	"delivery-backend/service/api/biz/dal/redis"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gin-gonic/gin"
)

func IPRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		warp := app.RespWarp{Context: c}
		key := c.ClientIP()
		remains, err := redis.AcquireBucket(key)
		if err != nil {
			klog.Error(err)
			warp.RespInnerError()
			return
		} else if remains <= 0 {
			klog.Error("remains token:", remains)
			warp.RespBadReq()
			return
		}
		klog.Debug("tokens remains: ", remains)
		c.Next()
	}
}
