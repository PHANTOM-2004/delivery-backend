
package app

import (
  "delivery-backend/internal/ecode"

	"github.com/gin-gonic/gin"
)

func Response(c *gin.Context, httpCode int, errCode ecode.Ecode, data any) {
	c.JSON(httpCode, gin.H{
		"code": httpCode,
		"msg":  ecode.StatusText(errCode),
		"data": data,
	})
}
