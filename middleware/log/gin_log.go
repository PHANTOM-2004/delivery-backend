package gin_log

import (
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Stop timer
		timestamp := time.Now()
		latency := timestamp.Sub(start)
		client_ip := c.ClientIP()
		method := c.Request.Method
		status_code := c.Writer.Status()
		uri := c.Request.RequestURI

		log.Infof(
			"[GIN-handle] %3d | %6s | %10v | %15s | %s",
			status_code,
			method,
			latency,
			client_ip,
			uri,
		)
	}
}
