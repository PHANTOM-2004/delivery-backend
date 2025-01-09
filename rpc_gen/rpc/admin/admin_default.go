package admin

import (
	"context"
	admin "delivery-backend/rpc_gen/kitex_gen/user/admin"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
)

func AdminRegister(ctx context.Context, req *admin.AdminRegisterReq, callOptions ...callopt.Option) (resp *admin.AdminRegisterResp, err error) {
	resp, err = defaultClient.AdminRegister(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "AdminRegister call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func AdminLogin(ctx context.Context, req *admin.AdminLoginReq, callOptions ...callopt.Option) (resp *admin.AdminLoginResp, err error) {
	resp, err = defaultClient.AdminLogin(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "AdminLogin call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}
