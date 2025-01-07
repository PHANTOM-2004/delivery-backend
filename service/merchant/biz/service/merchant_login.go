package service

import (
	"context"
	merchant "delivery-backend/rpc_gen/kitex_gen/user/merchant"
)

type MerchantLoginService struct {
	ctx context.Context
} // NewMerchantLoginService new MerchantLoginService
func NewMerchantLoginService(ctx context.Context) *MerchantLoginService {
	return &MerchantLoginService{ctx: ctx}
}

// Run create note info
func (s *MerchantLoginService) Run(req *merchant.MerchantLoginReq) (resp *merchant.MerchantLoginResp, err error) {
	// Finish your business logic.

	return
}
