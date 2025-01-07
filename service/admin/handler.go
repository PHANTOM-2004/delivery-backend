package main

import (
	"context"
	admin "delivery-backend/rpc_gen/kitex_gen/user/admin"
	"delivery-backend/service/admin/biz/service"
)

// AdminServiceImpl implements the last service interface defined in the IDL.
type AdminServiceImpl struct{}

// AdminRegister implements the AdminServiceImpl interface.
func (s *AdminServiceImpl) AdminRegister(ctx context.Context, req *admin.AdminRegisterReq) (resp *admin.AdminRegisterResp, err error) {
	resp, err = service.NewAdminRegisterService(ctx).Run(req)

	return resp, err
}

// AdminLogin implements the AdminServiceImpl interface.
func (s *AdminServiceImpl) AdminLogin(ctx context.Context, req *admin.AdminLoginReq) (resp *admin.AdminLoginResp, err error) {
	resp, err = service.NewAdminLoginService(ctx).Run(req)

	return resp, err
}
