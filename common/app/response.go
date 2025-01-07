package app

import (
	"delivery-backend/common/ecode"
	"errors"
	"net/http"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gin-gonic/gin"
)

type RespWarp struct {
	Context *gin.Context
}

func (w *RespWarp) RespInnerError() {
	// klog.Warn(err)
	w.Context.JSON(http.StatusInternalServerError, gin.H{
		"ecode": ecode.ERROR,
		"msg":   ecode.StatusText(ecode.ERROR),
		"data":  nil,
	})
	w.Abort()
}

func (w *RespWarp) RespBadReq() {
	w.Context.JSON(http.StatusBadRequest, gin.H{
		"ecode": ecode.INVALID_PARAMS,
		"msg":   ecode.StatusText(ecode.INVALID_PARAMS),
		"data":  nil,
	})
	w.Abort()
}

func (w *RespWarp) RespSuccData(data map[string]any) {
	w.Context.JSON(http.StatusOK, gin.H{
		"ecode": ecode.SUCCESS,
		"msg":   ecode.StatusText(ecode.SUCCESS),
		"data":  data,
	})
}

func (w *RespWarp) RespSucc() {
	w.Context.JSON(http.StatusOK, gin.H{
		"ecode": ecode.SUCCESS,
		"msg":   ecode.StatusText(ecode.SUCCESS),
		"data":  nil,
	})
}

// should resp biz
func (w *RespWarp) RespRPCErr(err error) {
	var bizErr *kerrors.GRPCBizStatusError
	ok := errors.As(err, &bizErr)
	if ok {
		w.Context.JSON(
			http.StatusOK,
			gin.H{
				"ecode": bizErr.BizStatusCode(),
				"msg":   bizErr.BizMessage(),
				"data":  nil,
			},
		)
		w.Abort()
	} else {
		klog.Error(err)
		w.RespInnerError()
	}
}

func (w *RespWarp) Resp(httpCode int, errCode int, data any) {
	w.Context.JSON(httpCode, gin.H{
		"ecode": errCode,
		"msg":   ecode.StatusText(errCode),
		"data":  data,
	})
}

func (w *RespWarp) Abort() {
	w.Context.Abort()
}
