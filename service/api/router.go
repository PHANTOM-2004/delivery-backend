package main

import (
	"delivery-backend/common/app"
	"delivery-backend/service/api/biz/v1/middleware"
	"delivery-backend/service/api/biz/v1/user"
	"delivery-backend/service/api/conf"
	"net/http"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func initMerchantRouter(r *gin.RouterGroup) {
	session_store, err := redis.NewStore(
		// TODO:
		10,
		"tcp",
		conf.GetConf().Redis.Address,
		conf.GetConf().Redis.Password,
		[]byte("666"),
	)
	if err != nil {
		klog.Fatal(err)
	}
	session_store.Options(sessions.Options{
		Path: "/api/v1/merchant",
		// TODO:
		MaxAge:   20 * 60,
		Secure:   true, // 仅通过HTTPS传输Cookie
		HttpOnly: true, // 禁止通过JavaScript访问Cookie
		SameSite: http.SameSiteDefaultMode,
	})

	sessionHandler := sessions.Sessions("MerchantSession", session_store)

	merchant := r.Group(
		"/merchant",
		sessionHandler,
	)
	// TODO:
	merchant.POST("/login", user.MerchantLogin)

	merchant_session := merchant.Group(
		"/",
		middleware.MerchAuth(),
	)
	merchant_session.GET("/login-status",
		func(c *gin.Context) {
			warp := app.RespWarp{Context: c}
			warp.RespSucc()
		})
	merchant_session.POST("/logout", user.MerchantLogout)
}

func initAdminRouter(r *gin.RouterGroup) {
	session_store, err := redis.NewStore(
		10,
		"tcp",
		conf.GetConf().Redis.Address,
		conf.GetConf().Redis.Password,
		[]byte("666"),
	)
	if err != nil {
		klog.Fatal(err)
	}
	session_store.Options(sessions.Options{
		Path:     "/api/v1/admin",
		MaxAge:   20 * 60,
		Secure:   true, // 仅通过HTTPS传输Cookie
		HttpOnly: true, // 禁止通过JavaScript访问Cookie
		SameSite: http.SameSiteDefaultMode,
	})
	admin_session_handler := sessions.Sessions("AdminSession", session_store)

	admin_account := r.Group("/admin")
	admin_account.Use(admin_session_handler)
	admin_account.POST("/create",
		user.AdminRegister)
	admin_account.POST("/login",
		user.AdminLogin)

	admin_session := admin_account.Group(
		"/",
		middleware.AdminAuth(),
	)

	admin_session.GET("/login-status",
		func(c *gin.Context) {
			warp := app.RespWarp{Context: c}
			warp.RespSucc()
		})
	admin_session.POST("/logout", user.AdminLogout)
  admin_session.POST("/merchant/create", user.CreateMerch)
}
