package service

import (
	"context"
	merchant "delivery-backend/rpc_gen/kitex_gen/user/merchant"
	"delivery-backend/service/merchant/biz/dal/model"
	"delivery-backend/service/merchant/biz/dal/mysql"
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
	m := model.Merchant{}
	err = mysql.DB.Find(&m, model.Merchant{Account: req.Account}).Error
	resp = &merchant.MerchantLoginResp{
		UserId: uint32(m.ID),
	}
	return
}
