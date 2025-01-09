package service

import (
	"context"
	merchant "delivery-backend/rpc_gen/kitex_gen/user/merchant"
	"testing"
)

func TestMerchantRegister_Run(t *testing.T) {
	ctx := context.Background()
	s := NewMerchantRegisterService(ctx)
	// init req and assert value

	req := &merchant.MerchantRegisterReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test
}
