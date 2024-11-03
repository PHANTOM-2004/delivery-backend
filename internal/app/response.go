package app

import (
	"delivery-backend/internal/ecode"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResponseInternalError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"ecode": ecode.ERROR,
		"msg":   ecode.StatusText(ecode.ERROR),
		"data":  nil,
	})
}

func ResponseSuccess(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ecode": ecode.SUCCESS,
		"msg":   ecode.StatusText(ecode.SUCCESS),
		"data":  nil,
	})
}

func Response(c *gin.Context, httpCode int, errCode ecode.Ecode, data any) {
	c.JSON(httpCode, gin.H{
		"ecode": errCode,
		"msg":   ecode.StatusText(errCode),
		"data":  data,
	})
}
