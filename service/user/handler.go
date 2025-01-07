package main

import (
	"context"
	user "delivery-backend/rpc_gen/kitex_gen/user"
	"delivery-backend/service/user/biz/service"
)

// AdminServiceImpl implements the last service interface defined in the IDL.
type AdminServiceImpl struct{}

// AdminRegister implements the AdminServiceImpl interface.
func (s *AdminServiceImpl) AdminRegister(ctx context.Context, req *user.AdminRegisterReq) (resp *user.AdminRegisterResp, err error) {
	resp, err = service.NewAdminRegisterService(ctx).Run(req)

	return resp, err
}

// AdminLogin implements the AdminServiceImpl interface.
func (s *AdminServiceImpl) AdminLogin(ctx context.Context, req *user.AdminLoginReq) (resp *user.AdminLoginResp, err error) {
	resp, err = service.NewAdminLoginService(ctx).Run(req)

	return resp, err
}
