package routers

import (
	"delivery-backend/internal/setting"
	"delivery-backend/middleware/jwt"
	"delivery-backend/routers/api"
	v1 "delivery-backend/routers/api/v1"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

func InitRouter() *gin.Engine {
	defer log.Info("app router initialized")

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(setting.ServerSetting.RunMode)

	r.POST("/admin/login", api.AdminLogin)
	r.GET("/admin/auth", api.GetAuth)
  r.POST("/admin/create", api.AdminCreate)

	{
		// TODO: 管理员修改密码
		apiv1 := r.Group("/api/v1")

		{
			// admin group
			admin := apiv1.Group("/admin")

			// middleware JWT
			admin.Use(jwt.JWT())

			// middleware redis session
			store, err := redis.NewStore(
				setting.RedisSetting.MaxIdle,
				"tcp",
				setting.RedisSetting.Host,
				setting.RedisSetting.Password,
				[]byte(setting.RedisSetting.Secret),
			)
			if err != nil {
				log.Fatal(err)
			}

			admin.Use(sessions.Sessions("AdminSession", store))
			admin.GET("/change-password", v1.AdminChangePassword)
		}

	}

	return r
}
