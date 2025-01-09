package service

import (
	"context"
	"delivery-backend/common/ecode"
	"delivery-backend/common/util"
	merchant "delivery-backend/rpc_gen/kitex_gen/user/merchant"
	"delivery-backend/service/merchant/biz/dal/model"
	"delivery-backend/service/merchant/biz/dal/mysql"
	"os"

	"github.com/cloudwego/kitex/pkg/kerrors"
)

type MerchantLoginService struct {
	ctx context.Context
} // NewMerchantLoginService new MerchantLoginService
func NewMerchantLoginService(ctx context.Context) *MerchantLoginService {
	return &MerchantLoginService{ctx: ctx}
}

var salt = os.Getenv("MERCH_PWD_SALT")

// Run create note info
func (s *MerchantLoginService) Run(req *merchant.MerchantLoginReq) (resp *merchant.MerchantLoginResp, err error) {
	// Finish your business logic.
	m := model.Merchant{}
	err = mysql.DB.Find(&m, model.Merchant{Account: req.Account}).Error
	resp = &merchant.MerchantLoginResp{
		UserId: uint32(m.ID),
	}

	if util.Encrypt(req.Password, salt) != m.Password {
		code := ecode.ERROR_MERCHANT_INCORRECT_PWD
		err = kerrors.NewGRPCBizStatusError(int32(code), ecode.StatusText(code))
		return
	}

	return
}
