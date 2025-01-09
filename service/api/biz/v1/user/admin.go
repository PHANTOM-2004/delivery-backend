package user

import (
	"context"
	"delivery-backend/common/app"
	"delivery-backend/common/ecode"
	"delivery-backend/common/util"
	"delivery-backend/rpc_gen/kitex_gen/user/admin"
	"delivery-backend/rpc_gen/kitex_gen/user/merchant"
	"net/http"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type adminLoginReq struct {
	Account  string `form:"account" validate:"min=10,max=30"`
	Password string `form:"password" validate:"min=15,max=30"`
}

func AdminLogout(c *gin.Context) {
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

func AdminLogin(c *gin.Context) {
	var req adminLoginReq
	resp := app.RespWarp{Context: c}
	// https://gin-gonic.com/docs/examples/only-bind-query-string/
	err := c.Bind(&req)
	if err != nil {
		klog.Debug(err)
		resp.RespBadReq()
		return
	}
	err = app.ValidateStruct(&req)
	if err != nil {
		klog.Debug(err)
		resp.RespBadReq()
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	rpcResp, err := rpcClientAdmin.AdminLogin(
		ctx,
		&admin.AdminLoginReq{
			Account:  req.Account,
			Password: req.Password,
		},
	)
	if err != nil {
		resp.RespRPCErr(err)
		return
	}
	if rpcResp.UserId == 0 {
		resp.Resp(
			http.StatusUnauthorized,
			ecode.ERROR_ADMIN_NOT_FOUND,
			nil,
		)
		return
	}

	// set redis
	session := sessions.Default(c)
	session.Set("admin_id", rpcResp.UserId)
	err = session.Save()
	if err != nil {
		klog.Error(err)
		resp.RespInnerError()
		return
	}

	resp.RespSucc()
}

type adminRegisterReq struct {
	Account  string `form:"account" validate:"max=30"`
	Password string `form:"password" validate:"min=15,max=30"`
	Name     string `form:"admin_name" validate:"min=2,max=20"`
}

// NOTE:该接口只允许运维调用，需要验证创建管理员的唯一token.
func AdminRegister(c *gin.Context) {
	// TODO: limit ip
	var req adminRegisterReq
	resp := app.RespWarp{Context: c}
	// https://gin-gonic.com/docs/examples/only-bind-query-string/
	err := c.ShouldBindQuery(&req)
	if err != nil {
		klog.Debug(err, req)
		resp.RespBadReq()
		return
	}
	err = app.ValidateStruct(&req)
	if err != nil || len(req.Account) < 10 {
		klog.Debug(err, req)
		resp.RespBadReq()
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// 分发给rpc处理
	rpcResp, err := rpcClientAdmin.AdminRegister(
		ctx,
		&admin.AdminRegisterReq{
			Name:     req.Name,
			Password: req.Password,
			Account:  req.Account,
		},
	)
	// check kerrors
	if err != nil {
		resp.RespRPCErr(err)
		return
	}

	// 交给答复
	resp.RespSuccData(map[string]any{
		"account": rpcResp.Account,
	})
}

type createMerchReq struct {
	Name        string `form:"merchant_name" validate:"min=2,max=20"`
	Account     string `form:"account" validate:"max=30"`  // min=6
	Password    string `form:"password" validate:"max=30"` // min=8
	PhoneNumber string `form:"phone_number" validate:"required,e164"`
	Email       string `form:"email" validate:"required,email"`
}

func CreateMerch(c *gin.Context) {
	var err error
	var req createMerchReq
	resp := app.RespWarp{Context: c}
	err = c.Bind(&req)
	if err != nil {
		klog.Debug(err, req)
		resp.RespBadReq()
		return
	}
	err = app.ValidateStruct(&req)
	if err != nil || len(req.Account) < 6 {
		klog.Debug(err, req)
		resp.RespBadReq()
		return
	}
	if req.Password == "" {
		req.Password = util.RandString(12)
	}
	if req.Account == "" {
		req.Account = util.RandString(12)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	createdDone := make(chan struct{})
	var createdErr error
	go func() {
		_, createdErr = rpcClientMerch.MerchantRegister(
			ctx,
			&merchant.MerchantRegisterReq{
				Name:        req.Name,
				Password:    req.Password,
				Account:     req.Account,
				PhoneNumber: req.PhoneNumber,
			},
		)
		createdDone <- struct{}{}
	}()
	// TODO:邮件服务

	select {
	case <-createdDone:
		klog.Debug("created merchant done")
		if createdErr != nil {
			resp.RespRPCErr(createdErr)
			return
		}
	}

	// 交给答复
	resp.RespSuccData(map[string]any{
		"account":  req.Account,
		"password": req.Password,
	})
}
