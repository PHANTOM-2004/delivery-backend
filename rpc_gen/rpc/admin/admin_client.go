package admin

import (
	"context"
	admin "delivery-backend/rpc_gen/kitex_gen/user/admin"

	"delivery-backend/rpc_gen/kitex_gen/user/admin/adminservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

type RPCClient interface {
	KitexClient() adminservice.Client
	Service() string
	AdminRegister(ctx context.Context, Req *admin.AdminRegisterReq, callOptions ...callopt.Option) (r *admin.AdminRegisterResp, err error)
	AdminLogin(ctx context.Context, Req *admin.AdminLoginReq, callOptions ...callopt.Option) (r *admin.AdminLoginResp, err error)
}

func NewRPCClient(dstService string, opts ...client.Option) (RPCClient, error) {
	kitexClient, err := adminservice.NewClient(dstService, opts...)
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
	kitexClient adminservice.Client
}

func (c *clientImpl) Service() string {
	return c.service
}

func (c *clientImpl) KitexClient() adminservice.Client {
	return c.kitexClient
}

func (c *clientImpl) AdminRegister(ctx context.Context, Req *admin.AdminRegisterReq, callOptions ...callopt.Option) (r *admin.AdminRegisterResp, err error) {
	return c.kitexClient.AdminRegister(ctx, Req, callOptions...)
}

func (c *clientImpl) AdminLogin(ctx context.Context, Req *admin.AdminLoginReq, callOptions ...callopt.Option) (r *admin.AdminLoginResp, err error) {
	return c.kitexClient.AdminLogin(ctx, Req, callOptions...)
}
