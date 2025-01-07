package service

import (
	"context"
	merchant "delivery-backend/rpc_gen/kitex_gen/user/merchant"
)

type MerchantApplyService struct {
	ctx context.Context
} // NewMerchantApplyService new MerchantApplyService
func NewMerchantApplyService(ctx context.Context) *MerchantApplyService {
	return &MerchantApplyService{ctx: ctx}
}

// Run create note info
func (s *MerchantApplyService) Run(req *merchant.MerchantApplyReq) (resp *merchant.MerchantApplyResp, err error) {
	// Finish your business logic.

	return
}
