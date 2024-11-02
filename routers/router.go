package routers

import (
	"delivery-backend/internal/setting"
	"delivery-backend/middleware/jwt"
	"delivery-backend/routers/api"

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
	// session for admin usage,暂时不使用
	admin_session_handler := sessions.Sessions("AdminSession", store)
	log.Debug("Currently session not used", admin_session_handler)

	// admin group, for vite usage
	admin := r.Group("/admin")
	// admin api
	admin.POST("/create", api.AdminCreate)
	admin.DELETE("/delete", api.AdminDelete)
	admin.POST("/login", api.AdminLogin)
	admin.POST("/logout", api.AdminLogout)

	// merchant
	merchant := r.Group("/merchant")
	merchant.POST("/login")

	{
		apiv1 := r.Group("/api/v1")

		{

			// admin group
			admin_jwt := apiv1.Group("/admin/jwt")

			// apis
			{
				admin_jwt_ak := admin_jwt.Group("/")
				admin_jwt_ak.Use(jwt.JWTAK())
				admin_jwt_ak.PUT("/change-password", api.AdminChangePassword)
			}

			{
				admin_jwt_rk := admin_jwt.Group("/")
				admin_jwt_rk.Use(jwt.JWTRK())
				admin_jwt_rk.GET("/auth", api.AdminGetAuth)
				admin_jwt_rk.GET("/login-status", api.AdminLoginStatus)
			}
		}
	}

	return r
}
