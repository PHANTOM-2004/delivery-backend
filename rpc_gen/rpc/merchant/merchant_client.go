package merchant

import (
	"context"
	merchant "delivery-backend/rpc_gen/kitex_gen/user/merchant"

	"delivery-backend/rpc_gen/kitex_gen/user/merchant/merchantservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

type RPCClient interface {
	KitexClient() merchantservice.Client
	Service() string
	MerchantRegister(ctx context.Context, Req *merchant.MerchantRegisterReq, callOptions ...callopt.Option) (r *merchant.MerchantRegisterResp, err error)
	MerchantLogin(ctx context.Context, Req *merchant.MerchantLoginReq, callOptions ...callopt.Option) (r *merchant.MerchantLoginResp, err error)
}

func NewRPCClient(dstService string, opts ...client.Option) (RPCClient, error) {
	kitexClient, err := merchantservice.NewClient(dstService, opts...)
	if err != nil {
		return nil, err
	}
	cli := &clientImpl{
		service:     dstService,
		kitexClient: kitexClient,
	}

	return cli, nil
}

type clientImpl struct {
	service     string
	kitexClient merchantservice.Client
}

func (c *clientImpl) Service() string {
	return c.service
}

func (c *clientImpl) KitexClient() merchantservice.Client {
	return c.kitexClient
}

func (c *clientImpl) MerchantRegister(ctx context.Context, Req *merchant.MerchantRegisterReq, callOptions ...callopt.Option) (r *merchant.MerchantRegisterResp, err error) {
	return c.kitexClient.MerchantRegister(ctx, Req, callOptions...)
}

func (c *clientImpl) MerchantLogin(ctx context.Context, Req *merchant.MerchantLoginReq, callOptions ...callopt.Option) (r *merchant.MerchantLoginResp, err error) {
	return c.kitexClient.MerchantLogin(ctx, Req, callOptions...)
}
