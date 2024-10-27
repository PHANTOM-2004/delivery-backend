package routers

import (
	"delivery-backend/internal/setting"
	"delivery-backend/middleware/filter"
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
	admin.POST("/login", admin_session_handler, api.AdminLogin)
	admin.POST("/logout", admin_session_handler, api.AdminLogout)
	admin.GET("/auth", api.GetAuth)
	admin.POST("/create", api.AdminCreate)
	admin.DELETE("/delete", api.AdminDelete)

	{
		// admin group
		admin_session := admin.Group("/session")
		// middleware session
		admin_session.Use(admin_session_handler)
		// middleware allowed only when logged in
		admin_session.Use(filter.LoginFilter())
		// middleware JWT
		admin_session.Use(jwt.JWT())
		// middleware double validation
		admin_session.Use(filter.DoubleValidation())

		// apis
		admin_session.PUT("/change-password", v1.AdminChangePassword)
	}

	// apiv1 := r.Group("/api/v1")
	return r
}
