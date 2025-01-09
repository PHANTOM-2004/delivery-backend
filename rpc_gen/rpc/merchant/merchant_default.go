package merchant

import (
	"context"
	merchant "delivery-backend/rpc_gen/kitex_gen/user/merchant"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
)

func MerchantRegister(ctx context.Context, req *merchant.MerchantRegisterReq, callOptions ...callopt.Option) (resp *merchant.MerchantRegisterResp, err error) {
	resp, err = defaultClient.MerchantRegister(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "MerchantRegister call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func MerchantLogin(ctx context.Context, req *merchant.MerchantLoginReq, callOptions ...callopt.Option) (resp *merchant.MerchantLoginResp, err error) {
	resp, err = defaultClient.MerchantLogin(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "MerchantLogin call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}
