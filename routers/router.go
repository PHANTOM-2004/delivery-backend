package routers

import (
	"delivery-backend/internal/setting"
	"delivery-backend/middleware/jwt"
	filter "delivery-backend/middleware/login"
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
	// NOTE:注意更新文档

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
	// session for admin usage
	adminSession := sessions.Sessions("AdminSession", store)

	r.POST("/admin/login", adminSession, api.AdminLogin)
	r.POST("/admin/logout", adminSession, api.AdminLogout)

	r.GET("/admin/auth", api.GetAuth)
	r.POST("/admin/create", api.AdminCreate)
	r.DELETE("/admin/delete", api.AdminDelete)

	{
		apiv1 := r.Group("/api/v1")

		{
			// admin group
			admin := apiv1.Group("/admin")
			// middleware session
			admin.Use(adminSession)
			// middleware allowed only when logged in
			admin.Use(filter.LoginFilter())
			// middleware JWT
			admin.Use(jwt.JWT())

			// apis
			admin.PUT("/change-password", v1.AdminChangePassword)
		}

	}

	return r
}
