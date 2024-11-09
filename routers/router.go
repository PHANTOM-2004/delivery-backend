package routers

import (
	"delivery-backend/internal/setting"
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
	r.MaxMultipartMemory = int64(setting.AppSetting.MaxImageSize << 20) // MiB

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
			// middleware: 过滤在黑名单中的商家
			merchant_jwt_ak.Use(ak_hanlder,
				merchant_service.MerchantBlacklistFilter())

			merchant_jwt_ak.PUT("/change-password",
				api.MerchantChangePassword)

			merchant_jwt_ak.GET("/restaurants",
				v1.GetRestaurants)

			merchant_jwt_ak.POST("/restaurant/create", v1.CreateRestaurant)

			{
				// NOTE: 鉴权，商家必须对自己的restaurant操作
				merchant_restaurant := merchant_jwt_ak.Group("/restaurant/:restaurant_id")

				merchant_restaurant.Use(
					merchant_service.RestaurantAuth(),
				)

				merchant_restaurant.PUT(
					"/update",
					v1.UpdateRestaurant,
				)

				merchant_restaurant.PUT(
					"/status/:status",
					v1.SetRestaurantStatus,
				)

				merchant_restaurant.GET(
					"/status",
					v1.GetRestaurantStatus,
				)

				merchant_restaurant.GET(
					"/categories",
					v1.GetCategories,
				)

				merchant_restaurant.PUT(
					"/category/:category_id/update",
					v1.UpdateCategory,
				)

				merchant_restaurant.POST(
					"/category/create",
					v1.CreateCategory,
				)

				merchant_restaurant.PUT(
					"/category/:category_id/dish/:dish_id/update",
					merchant_service.CategoryAuth(),
					v1.UpdateDish,
				)

				merchant_restaurant.POST(
					"/category/:category_id/dish/create",
					merchant_service.CategoryAuth(),
					v1.CreateDish,
				)

				// NOTE: license的图片静态文件路由
				dish_image_path := setting.AppSetting.DishImageStorePath
				log.Infof("Serving Static File: [%s]", dish_image_path)
				merchant_jwt_ak.Static("/dish", dish_image_path)
			}
		}

		{
			merchant_jwt_rk := merchant_jwt.Group("/")
			rk_hanlder := jwt.JWTRK(
				merchant_service.TokenInBlacklist,
				merchant_service.AuthRefreshToken,
			)
			// middleware: 过滤在黑名单中的商家
			merchant_jwt_rk.Use(rk_hanlder,
				merchant_service.MerchantBlacklistFilter())
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
				api.CreateMerchant)
			admin_jwt_ak.POST("/merchant-delete",
				api.DeleteMerchant)
			admin_jwt_ak.GET("/merchant-application/:page",
				v1.GetMerchantApplication)
			admin_jwt_ak.PUT("/merchant-application/:application_id/approve",
				v1.ApproveMerchantApplication)
			admin_jwt_ak.PUT("/merchant-application/:application_id/disapprove",
				v1.DisapproveMerchantApplication)

			// NOTE: license的图片静态文件路由
			license_path := setting.AppSetting.LicenseStorePath
			log.Infof("Serving Static File: [%s]", license_path)
			admin_jwt_ak.Static("/merchant-application/license", license_path)
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
