package user

import (
	"context"
	"delivery-backend/common/app"
	"delivery-backend/rpc_gen/kitex_gen/user/admin"
	"delivery-backend/rpc_gen/kitex_gen/user/admin/adminservice"
	"delivery-backend/service/api/conf"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var rpcClient adminservice.Client

// 初始化etcd client
func init() {
	r, err := etcd.NewEtcdResolver(
		conf.GetConf().Registry.RegistryAddress,
	)
	if err != nil {
		klog.Fatal(err)
	}
	rpcClient, err = adminservice.NewClient("admin", client.WithResolver(r),
		// NOTE: 注意选择正确协议, 如果不选择正确协议或者协议不写, 我们就接受不到error
		client.WithMetaHandler(transmeta.ClientHTTP2Handler),
		client.WithTransportProtocol(transport.GRPC),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: "user.admin",
		}),
	)
	if err != nil {
		klog.Fatal(err)
	}
	klog.Info("RPC Client Init Done")
}

type adminLoginReq struct {
	Account  string `form:"account" validate:"min=10,max=30"`
	Password string `form:"password" validate:"min=15,max=30"`
}

func AdminLogout(c *gin.Context) {
	resp := app.RespWarp{Context: c}
	session := sessions.Default(c)
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	err := session.Save()
	if err != nil {
		// redis err
		resp.RespInnerError()
		klog.Error(err)
		return
	}
	resp.RespSucc()
}

func AdminLogin(c *gin.Context) {
	var req adminLoginReq
	resp := app.RespWarp{Context: c}
	// https://gin-gonic.com/docs/examples/only-bind-query-string/
	err := c.Bind(&req)
	if err != nil {
		klog.Debug(err)
		resp.RespBadReq()
		return
	}
	err = app.ValidateStruct(&req)
	if err != nil {
		klog.Debug(err)
		resp.RespBadReq()
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	rpcResp, err := rpcClient.AdminLogin(
		ctx,
		&admin.AdminLoginReq{
			Account:  req.Account,
			Password: req.Password,
		},
	)
	if err != nil {
		resp.RespRPCErr(err)
		return
	}

	// set redis
	session := sessions.Default(c)
	session.Set("admin_id", rpcResp.UserId)
	err = session.Save()
	if err != nil {
		klog.Error(err)
		resp.RespInnerError()
		return
	}

	resp.RespSucc()
}

type adminRegisterReq struct {
	Account  string `form:"account" validate:"max=30"`
	Password string `form:"password" validate:"min=15,max=30"`
	Name     string `form:"admin_name" validate:"min=2,max=20"`
}

// NOTE:该接口只允许运维调用，需要验证创建管理员的唯一token.
func AdminRegister(c *gin.Context) {
	// TODO: limit ip
	var req adminRegisterReq
	resp := app.RespWarp{Context: c}
	// https://gin-gonic.com/docs/examples/only-bind-query-string/
	err := c.ShouldBindQuery(&req)
	if err != nil {
		resp.RespBadReq()
		return
	}
	err = app.ValidateStruct(&req)
	if err != nil || len(req.Account) < 10 {
		resp.RespBadReq()
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// 分发给rpc处理
	rpcResp, err := rpcClient.AdminRegister(
		ctx,
		&admin.AdminRegisterReq{
			Name:     req.Name,
			Password: req.Password,
			Account:  req.Account,
		},
	)
	// check kerrors
	if err != nil {
		resp.RespRPCErr(err)
		return
	}

	// 交给答复
	resp.RespSuccData(map[string]any{
		"account": rpcResp.Account,
	})
}
