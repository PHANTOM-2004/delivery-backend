package main

import (
	"delivery-backend/service/api/biz/v1/user"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func initMerchantRouter(r *gin.RouterGroup) {
	session_store, err := redis.NewStore(
		// TODO:
		10,
		"tcp",
		"127.0.0.1:6379",
		"",
	)
	if err != nil {
		log.Fatal(err)
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
}
