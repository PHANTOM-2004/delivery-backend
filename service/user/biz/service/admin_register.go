package service

import (
	"context"
	user "delivery-backend/rpc_gen/kitex_gen/user"
)

type AdminRegisterService struct {
	ctx context.Context
} // NewAdminRegisterService new AdminRegisterService
func NewAdminRegisterService(ctx context.Context) *AdminRegisterService {
	return &AdminRegisterService{ctx: ctx}
}

// Run create note info
func (s *AdminRegisterService) Run(req *user.AdminRegisterReq) (resp *user.AdminRegisterResp, err error) {
	// Finish your business logic.

	return
}
