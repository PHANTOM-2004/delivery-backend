package service

import (
	"context"
	user "delivery-backend/rpc_gen/kitex_gen/user"
)

type AdminLoginService struct {
	ctx context.Context
} // NewAdminLoginService new AdminLoginService
func NewAdminLoginService(ctx context.Context) *AdminLoginService {
	return &AdminLoginService{ctx: ctx}
}

// Run create note info
func (s *AdminLoginService) Run(req *user.AdminLoginReq) (resp *user.AdminLoginResp, err error) {
	// Finish your business logic.

	return
}
