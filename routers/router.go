package routers

import (
	"delivery-backend/internal/setting"
	"delivery-backend/routers/api"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

func InitRouter() *gin.Engine {
	defer log.Info("app router initialized")

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(setting.ServerSetting.RunMode)

	r.GET("/login/admin/auth", api.ValidateAccount)

	// admin := r.Group("/admin")

	// apiv1 := r.Group("/api/v1")
	// TODO: JWT 鉴权

	return r
}
