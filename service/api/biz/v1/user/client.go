package user

import (
	"delivery-backend/rpc_gen/kitex_gen/user/admin/adminservice"
	"delivery-backend/rpc_gen/kitex_gen/user/merchant/merchantservice"
	"delivery-backend/service/api/conf"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var (
	rpcClientAdmin adminservice.Client
	rpcClientMerch merchantservice.Client
)

// 初始化etcd client
func init() {
	r, err := etcd.NewEtcdResolver(
		conf.GetConf().Registry.RegistryAddress,
	)
	if err != nil {
		klog.Fatal(err)
	}
	rpcClientAdmin, err = adminservice.NewClient(
		"admin", client.WithResolver(r),
		// NOTE: 注意选择正确协议, 如果不选择正确协议或者协议不写, 我们就接受不到error
		client.WithMetaHandler(transmeta.ClientHTTP2Handler),
		client.WithTransportProtocol(transport.GRPC),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: "api.admin",
		}),
	)
	if err != nil {
		klog.Fatal(err)
	}
	klog.Info("Admin RPC Client Init Done")
}

func init() {
	r, err := etcd.NewEtcdResolver(
		conf.GetConf().Registry.RegistryAddress,
	)
	if err != nil {
		klog.Fatal(err)
	}
	rpcClientMerch, err = merchantservice.NewClient(
		"merchant", client.WithResolver(r),
		// NOTE: 注意选择正确协议, 如果不选择正确协议或者协议不写, 我们就接受不到error
		client.WithMetaHandler(transmeta.ClientHTTP2Handler),
		client.WithTransportProtocol(transport.GRPC),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: "api.merch",
		}),
	)
	if err != nil {
		klog.Fatal(err)
	}
	klog.Info("Merchant RPC Client Init Done")
}
