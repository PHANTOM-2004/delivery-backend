package app

import (
	"delivery-backend/internal/ecode"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func ResponseInternalError(c *gin.Context, err error) {
	log.Warn(err)
	c.JSON(http.StatusInternalServerError, gin.H{
		"ecode": ecode.ERROR,
		"msg":   ecode.StatusText(ecode.ERROR),
		"data":  nil,
	})
	c.Abort()
}

func ResponseInvalidParams(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"ecode": ecode.INVALID_PARAMS,
		"msg":   ecode.StatusText(ecode.INVALID_PARAMS),
		"data":  nil,
	})
	c.Abort()
}

func ResponseSuccessWithData(c *gin.Context, data map[string]any) {
	c.JSON(http.StatusOK, gin.H{
		"ecode": ecode.SUCCESS,
		"msg":   ecode.StatusText(ecode.SUCCESS),
		"data":  data,
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
