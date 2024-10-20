package routers

import (
	"delivery-backend/internal/setting"
	"delivery-backend/middleware/jwt"
	"delivery-backend/routers/api"
	v1 "delivery-backend/routers/api/v1"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

func InitRouter() *gin.Engine {
	defer log.Info("app router initialized")

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(setting.ServerSetting.RunMode)

	r.GET("/admin/login", api.AdminLogin)
	r.GET("/admin/auth", api.GetAuth)

	// TODO: 管理员修改密码, after JWT鉴权

	admin := r.Group("/api/v1/admin")
	admin.Use(jwt.JWT())
	admin.GET("/change-password", v1.AdminChangePassword)

	// admin := r.Group("/admin")

	// apiv1 := r.Group("/api/v1")
	// TODO: JWT 鉴权

	return r
}
