package routers

import (
	"delivery-backend/internal/setting"
	"delivery-backend/middleware/filter"
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
	// session for admin usage
	admin_session_handler := sessions.Sessions("AdminSession", store)
	// admin group, for vite usage
	admin := r.Group("/admin")
	// admin api
	admin.GET("/auth", api.GetAuth)
	admin.POST("/create", api.AdminCreate)
	admin.DELETE("/delete", api.AdminDelete)
	admin.POST("/login", admin_session_handler, api.AdminLogin)

	{
		// NOTE: 对于有状态服务，采取session已经是比较好的选择
		// 因此不再使用JWT鉴权

		// admin group
		admin_session := admin.Group("/session")
		// middleware session
		admin_session.Use(admin_session_handler)
		// middleware allowed only when logged in
		admin_session.Use(filter.LoginFilter())

		// apis
		admin_session.POST("/logout", api.AdminLogout)
		admin_session.GET("/login-status", api.AdminLoginStatus)
		admin_session.PUT("/change-password", api.AdminChangePassword)
	}

	// apiv1 := r.Group("/api/v1")
	return r
}
