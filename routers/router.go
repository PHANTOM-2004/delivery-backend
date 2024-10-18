package routers

import (
	"delivery-backend/internal/setting"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

func InitRouter() *gin.Engine {
	defer log.Info("router initialized")

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(setting.ServerSetting.RunMode)

	// apiv1 := r.Group("/api/v1")
	// TODO: JWT 鉴权

	return r
}
