package main

import (
	"delivery-backend/service/api/biz/v1/user"
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
		"127.0.0.1:6379",
		"",
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

	merchant_session := sessions.Sessions("MerchantSession", session_store)

	merchant := r.Group("/merchant")
	merchant.Use(merchant_session)

	// TODO:
	merchant.POST("/login", user.MerchantLogin)
	merchant.POST("/logout", user.MerchantLogout)
}

func initAdminRouter(r *gin.RouterGroup) {
	session_store, err := redis.NewStore(
		10,
		"tcp",
		"127.0.0.1:6379",
		"",
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

	admin := r.Group("/admin")
	admin.POST("/create", user.AdminRegister)

	admin_session := admin.Group("/", admin_session_handler)
	admin_session.POST("/login", user.AdminLogin)
	admin_session.POST("/logout", user.AdminLogout)
}
