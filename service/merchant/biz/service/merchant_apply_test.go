package service

import (
	"context"
	"testing"
	merchant "delivery-backend/rpc_gen/kitex_gen/user/merchant"
)

func TestMerchantApply_Run(t *testing.T) {
	ctx := context.Background()
	s := NewMerchantApplyService(ctx)
	// init req and assert value

	req := &merchant.MerchantApplyReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
