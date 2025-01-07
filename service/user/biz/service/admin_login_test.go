package service

import (
	"context"
	"testing"
	user "delivery-backend/rpc_gen/kitex_gen/user"
)

func TestAdminLogin_Run(t *testing.T) {
	ctx := context.Background()
	s := NewAdminLoginService(ctx)
	// init req and assert value

	req := &user.AdminLoginReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
