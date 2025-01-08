package user

import (
	"context"
	"delivery-backend/common/app"
	"delivery-backend/common/ecode"
	"delivery-backend/rpc_gen/kitex_gen/user/merchant"
	"net/http"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type merchLoginReq struct {
	Account  string `form:"account" validate:"min=6,max=30"`
	Password string `form:"password" validate:"min=8,max=30"`
}

// 认证成功时会在该函数中设置
// c.Set("merchant_id", id)
func MerchantLogin(c *gin.Context) {
	resp := app.RespWarp{Context: c}
	var mreq merchLoginReq
	err := c.Bind(&mreq)
	if err != nil {
		klog.Debug(err)
		resp.RespBadReq()
		return
	}
	err = app.ValidateStruct(&resp)
	if err != nil {
		klog.Debug(err)
		resp.RespBadReq()
		return
	}
	// 交给rpc客户端
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*10,
	)
	defer cancel()
	rpcResp, err := rpcClientMerch.MerchantLogin(ctx,
		&merchant.MerchantLoginReq{
			Account:  mreq.Account,
			Password: mreq.Password,
		})
	if err != nil {
		resp.RespRPCErr(err)
		return
	}

	// 商家不存在
	if rpcResp.UserId == 0 {
		code := ecode.ERROR_MERCHANT_NON_FOUND
		resp.Resp(http.StatusUnauthorized,
			code, nil)
		return
	}

	session := sessions.Default(c)
	session.Set("merchant_id", rpcResp.UserId)
	err = session.Save()
	if err != nil {
		klog.Error(err)
		resp.RespInnerError()
		return
	}

	resp.RespSucc()
}

func MerchantLogout(c *gin.Context) {
	resp := app.RespWarp{Context: c}
	session := sessions.Default(c)
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	err := session.Save()
	if err != nil {
		// redis err
		resp.RespInnerError()
		klog.Error(err)
		return
	}
	resp.RespSucc()
}
