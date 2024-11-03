package routers

import (
	"delivery-backend/internal/setting"
	"delivery-backend/middleware/filter"
	"delivery-backend/middleware/jwt"
	"delivery-backend/routers/api"
	v1 "delivery-backend/routers/api/v1"
	"delivery-backend/service/admin_service"
	"delivery-backend/service/merchant_service"

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
	r.MaxMultipartMemory = int64(setting.AppSetting.LicensePhotoSize << 20) // MiB

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

	// admin group
	admin := r.Group("/admin")
	admin.POST("/create", api.AdminCreate)
	admin.DELETE("/delete", api.AdminDelete)
	admin.POST("/login", api.AdminLogin)
	admin.POST("/logout", api.AdminLogout)
	// merchant
	merchant := r.Group("/merchant")
	merchant.POST("/login", api.MerchantLogin)
	merchant.POST("/logout", api.MerchantLogout)

	apiv1 := r.Group("/api/v1")

	{
		// customer group
		// TODO: 身份校验
		customer := apiv1.Group("/customer")
		customer.POST("/merchant-application", v1.MerchantApply)
	}

	{
		// merchant group
		merchant_jwt := apiv1.Group("/merchant/jwt")
		{
			merchant_jwt_ak := merchant_jwt.Group("/")
			ak_hanlder := jwt.JWTAK(
				merchant_service.TokenInBlacklist,
				merchant_service.AuthAccessToken,
			)
			merchant_jwt_ak.Use(ak_hanlder, filter.MerchantBlacklistFilter())

			merchant_jwt_ak.PUT("/change-password",
				api.MerchantChangePassword)
		}

		{
			merchant_jwt_rk := merchant_jwt.Group("/")
			rk_hanlder := jwt.JWTRK(
				merchant_service.TokenInBlacklist,
				merchant_service.AuthRefreshToken,
			)
			merchant_jwt_rk.Use(rk_hanlder, filter.MerchantBlacklistFilter())
			merchant_jwt_rk.GET("/auth",
				api.MerchantGetAuth)
			merchant_jwt_rk.GET("/login-status",
				api.MerchantLoginStatus)
		}
	}

	{
		// admin group
		admin_jwt := apiv1.Group("/admin/jwt")
		{
			admin_jwt_ak := admin_jwt.Group("/")
			ak_hanlder := jwt.JWTAK(
				admin_service.TokenInBlacklist,
				admin_service.AuthAccessToken,
			)
			admin_jwt_ak.Use(ak_hanlder)
			admin_jwt_ak.PUT("/change-password", api.AdminChangePassword)
			admin_jwt_ak.POST("/merchant-create",
				api.MerchantCreate)
			admin_jwt_ak.POST("/merchant-delete",
				api.MerchantDelete)
			admin_jwt_ak.GET("/merchant-application/:page",
				v1.GetMerchantApplication)
			admin_jwt_ak.POST("/merchant-application/:id/approve",
				v1.ApproveMerchantApplication)
			admin_jwt_ak.POST("/merchant-application/:id/disapprove",
				v1.DisapproveMerchantApplication)
		}

		{
			admin_jwt_rk := admin_jwt.Group("/")
			rk_hanlder := jwt.JWTRK(
				admin_service.TokenInBlacklist,
				admin_service.AuthRefreshToken,
			)
			admin_jwt_rk.Use(rk_hanlder)
			admin_jwt_rk.GET("/auth",
				api.AdminGetAuth)
			admin_jwt_rk.GET("/login-status",
				api.AdminLoginStatus)
		}
	}

	return r
}
