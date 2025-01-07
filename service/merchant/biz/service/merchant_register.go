package service

import (
	"context"
	merchant "delivery-backend/rpc_gen/kitex_gen/user/merchant"
)

type MerchantRegisterService struct {
	ctx context.Context
} // NewMerchantRegisterService new MerchantRegisterService
func NewMerchantRegisterService(ctx context.Context) *MerchantRegisterService {
	return &MerchantRegisterService{ctx: ctx}
}

// Run create note info
func (s *MerchantRegisterService) Run(req *merchant.MerchantRegisterReq) (resp *merchant.MerchantRegisterResp, err error) {
	// Finish your business logic.

	return
}
