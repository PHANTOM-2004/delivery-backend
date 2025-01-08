package service

import (
	"context"
	"delivery-backend/common/ecode"
	merchant "delivery-backend/rpc_gen/kitex_gen/user/merchant"
	"delivery-backend/service/merchant/biz/dal/model"
	"delivery-backend/service/merchant/biz/dal/mysql"

	"github.com/cloudwego/kitex/pkg/kerrors"
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
	m := model.Merchant{
		MerchantName: req.Name,
		Account:      req.Account,
		Password:     req.Password,
		PhoneNumber:  req.PhoneNumber,
	}
	created, err := mysql.CreateMerchant(&m)
	if err != nil {
		return
	}

	if !created {
		code := ecode.ERROR_MERCHANT_ACCOUNT_EXIST
		err = kerrors.NewGRPCBizStatusError(int32(code), ecode.StatusText(code))
	}

	resp = &merchant.MerchantRegisterResp{
		UserId:  m.ID,
		Account: m.Account,
	}

	return
}
