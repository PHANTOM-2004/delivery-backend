package service

import (
	"context"
	"delivery-backend/common/ecode"
	"delivery-backend/common/util"
	admin "delivery-backend/rpc_gen/kitex_gen/user/admin"
	"delivery-backend/service/admin/biz/dal/model"
	"delivery-backend/service/admin/biz/dal/mysql"

	"github.com/cloudwego/kitex/pkg/kerrors"
)

type AdminRegisterService struct {
	ctx context.Context
} // NewAdminRegisterService new AdminRegisterService
func NewAdminRegisterService(ctx context.Context) *AdminRegisterService {
	return &AdminRegisterService{ctx: ctx}
}

// Run create note info
func (s *AdminRegisterService) Run(req *admin.AdminRegisterReq) (resp *admin.AdminRegisterResp, err error) {
	// Finish your business logic.
	req.Password = util.Encrypt(req.Password, "TODO:")
	// 如果Account为空, 生成随机12位的Account
	if req.Account == "" {
		req.Account = util.RandString(12)
	}
	a := model.Admin{
		AdminName: req.Name,
		Account:   req.Account,
		Password:  req.Password,
	}

	res := mysql.DB.Where(model.Admin{Account: req.Account}).
		FirstOrCreate(&a)
	if res.RowsAffected == 0 {
		// 封装error
		code := ecode.ERROR_ADMIN_ACCOUNT_EXIST
		err = kerrors.NewGRPCBizStatusError(
			int32(code),
			ecode.StatusText(code),
		)
		return
	} else if res.Error != nil {
		return
	}

	resp = &admin.AdminRegisterResp{
		AdminId: a.ID,
		Account: a.Account,
	}
	return
}
