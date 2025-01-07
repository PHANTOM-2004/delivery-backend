package service

import (
	"context"
	"testing"
	user "delivery-backend/rpc_gen/kitex_gen/user"
)

func TestAdminRegister_Run(t *testing.T) {
	ctx := context.Background()
	s := NewAdminRegisterService(ctx)
	// init req and assert value

	req := &user.AdminRegisterReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
