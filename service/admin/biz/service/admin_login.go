package service

import (
	"context"
	admin "delivery-backend/rpc_gen/kitex_gen/user/admin"
	"delivery-backend/service/admin/biz/dal/model"
	"delivery-backend/service/admin/biz/dal/mysql"
)

type AdminLoginService struct {
	ctx context.Context
} // NewAdminLoginService new AdminLoginService
func NewAdminLoginService(ctx context.Context) *AdminLoginService {
	return &AdminLoginService{ctx: ctx}
}

// Run create note info
func (s *AdminLoginService) Run(req *admin.AdminLoginReq) (resp *admin.AdminLoginResp, err error) {
	// Finish your business logic.
	a := model.Admin{}
	err = mysql.DB.Find(&a, model.Admin{Account: req.Account}).Error
	resp = &admin.AdminLoginResp{
		UserId: uint32(a.ID),
	}
	return
}
