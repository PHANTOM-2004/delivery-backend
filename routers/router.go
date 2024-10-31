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
	admin.GET("/auth", api.GetAuth)
	admin.POST("/create", api.AdminCreate)
	admin.DELETE("/delete", api.AdminDelete)
	admin.POST("/login", api.AdminLogin)

	{
		// admin group
		admin_jwt := admin.Group("/jwt")
		// jwt access_token
		admin_jwt.Use(jwt.JWT())

		// apis
		admin_jwt.POST("/logout", api.AdminLogout)
		admin_jwt.GET("/login-status", api.AdminLoginStatus)
		admin_jwt.PUT("/change-password", api.AdminChangePassword)
	}

	// apiv1 := r.Group("/api/v1")
	return r
}
