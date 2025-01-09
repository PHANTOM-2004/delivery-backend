package service

import (
	"context"
	"delivery-backend/common/ecode"
	"delivery-backend/common/util"
	admin "delivery-backend/rpc_gen/kitex_gen/user/admin"
	"delivery-backend/service/admin/biz/dal/model"
	"delivery-backend/service/admin/biz/dal/mysql"
	"os"

	"github.com/cloudwego/kitex/pkg/kerrors"
)

type AdminLoginService struct {
	ctx context.Context
} // NewAdminLoginService new AdminLoginService
func NewAdminLoginService(ctx context.Context) *AdminLoginService {
	return &AdminLoginService{ctx: ctx}
}

var salt = os.Getenv("ADMIN_PWD_SALT")

// Run create note info
func (s *AdminLoginService) Run(req *admin.AdminLoginReq) (resp *admin.AdminLoginResp, err error) {
	// Finish your business logic.
	a := model.Admin{}
	err = mysql.DB.Find(&a, model.Admin{Account: req.Account}).Error
	resp = &admin.AdminLoginResp{
		UserId: uint32(a.ID),
	}
	if util.Encrypt(req.Password, salt) != a.Password {
		code := ecode.ERROR_ADMIN_INCORRECT_PWD
		err = kerrors.NewGRPCBizStatusError(int32(code), ecode.StatusText(code))
		return
	}
	return
}
