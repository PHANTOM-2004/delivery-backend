package main

import (
	"context"
	merchant "delivery-backend/rpc_gen/kitex_gen/user/merchant"
	"delivery-backend/service/merchant/biz/service"
)

// MerchantServiceImpl implements the last service interface defined in the IDL.
type MerchantServiceImpl struct{}

// MerchantRegister implements the MerchantServiceImpl interface.
func (s *MerchantServiceImpl) MerchantRegister(ctx context.Context, req *merchant.MerchantRegisterReq) (resp *merchant.MerchantRegisterResp, err error) {
	resp, err = service.NewMerchantRegisterService(ctx).Run(req)

	return resp, err
}

// MerchantLogin implements the MerchantServiceImpl interface.
func (s *MerchantServiceImpl) MerchantLogin(ctx context.Context, req *merchant.MerchantLoginReq) (resp *merchant.MerchantLoginResp, err error) {
	resp, err = service.NewMerchantLoginService(ctx).Run(req)

	return resp, err
}

// MerchantApply implements the MerchantServiceImpl interface.
func (s *MerchantServiceImpl) MerchantApply(ctx context.Context, req *merchant.MerchantApplyReq) (resp *merchant.MerchantApplyResp, err error) {
	resp, err = service.NewMerchantApplyService(ctx).Run(req)

	return resp, err
}
