package routers

import (
	"delivery-backend/internal/setting"
	"delivery-backend/middleware/auth"
	"delivery-backend/middleware/wechat"
	"delivery-backend/routers/api"
	v1 "delivery-backend/routers/api/v1"
	"delivery-backend/service/merchant_service"
	"net/http"

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

	apiv1 := r.Group("/api/v1")

	wx := apiv1.Group("/wx")
	{
		// 微信group
		wx.POST("/login", v1.WXLogin)
	}

	/////////////////////////////////////////////////////
	/////////////////////////////////////////////////////
	/////////////////////////////////////////////////////
	/////////////////////////////////////////////////////

	{
		// customer group
		customer := wx.Group("/customer")
		// NOTE: user session 可以暂时不使用，方便测试
		customer.Use(wechat.WXsession())
		customer.POST("/info", v1.WXUploadUserInfo)
		customer.POST("/merchant-application", v1.MerchantApply)
		customer.GET("/restaurant/:restaurant_id/categories/dishes",
			v1.GetRestaurantCategoryDish)
		customer.GET("/restaurants", v1.WXGetRestaurants)
		customer.GET("/restaurant/:restaurant_id/dishes/top",
			v1.WXGetTopDishes)
		customer.POST("/cart/restaurant/:restaurant_id",
			v1.WXUpdateCart)
		customer.GET("/cart/restaurant/:restaurant_id",
			v1.WXGetCart)
		customer.GET("/addressbook", v1.GetAddressBook)
		customer.POST("/addressbook", v1.CreateAddressBook)
		customer.PUT("/addressbook/:address_book_id", v1.UpdateAddressBook)
		customer.PUT("/addressbook/:address_book_id/default", v1.SetDefaultAddressBook)
		customer.DELETE("/addressbook/:address_book_id", v1.DeleteAddressBook)
		customer.GET("/orders", v1.GetCustomerOrders)
		customer.POST("/order/restaurant/:restaurant_id", v1.CreateOrder)
		customer.POST("/order/:order_id/cancel", v1.CancelOrder)
		customer.POST("/comment/image", v1.WXUploadCommentImage)
    customer.POST("/comment/restaurant/:restaurant_id", v1.WXCreateComment)

		// 文件服务
		comment_image_path := setting.WechatSetting.CommentImageStorePath
		log.Infof("Serving Static File: [%s]", comment_image_path)
		customer.Static("/comment/image", comment_image_path)
	}

	////////////////////////////////////////////////////
	/////////////////////////////////////////////////////
	/////////////////////////////////////////////////////
	/////////////////////////////////////////////////////
	//
	{
		// middleware redis session
		merchant_session_store_v1, err := redis.NewStore(
			setting.RedisSetting.MaxIdle,
			"tcp",
			setting.RedisSetting.Host,
			setting.RedisSetting.Password,
			[]byte(setting.RedisSetting.Secret),
		)
		if err != nil {
			log.Fatal(err)
		}
		merchant_session_store_v1.Options(sessions.Options{
			Path:     "/api/v1/merchant",
			MaxAge:   setting.AppSetting.MerchantAliveMinute * 60,
			Secure:   true, // 仅通过HTTPS传输Cookie
			HttpOnly: true, // 禁止通过JavaScript访问Cookie
			SameSite: http.SameSiteDefaultMode,
		})
		merchant_session_handler := sessions.Sessions("MerchantSession", merchant_session_store_v1)

		merchantv1 := apiv1.Group("/merchant")
		merchantv1.Use(merchant_session_handler)
		merchantv1.POST("/login", api.MerchantLogin)
		merchantv1.POST("/logout", auth.MerchantAuth(), api.MerchantLogout)

		{
			// TODO:商家黑名单过滤

			// merchant group
			merchant_session := apiv1.Group("/merchant")
			merchant_session.Use(merchant_session_handler)
			merchant_session.Use(merchant_service.MerchantBlacklistFilter())

			///////////////商家账户相关
			merchant_session.GET("/login-status",
				v1.MerchantLoginStatus)
			merchant_session.PUT("/change-password",
				v1.MerchantChangePassword)
			merchant_session.GET("/info",
				v1.GetMerchantInfo)

			// NOTE: 鉴权，商家必须对自己的restaurant操作
			merchant_restaurant := merchant_session.Group("/restaurant/:restaurant_id")
			merchant_restaurant.Use(merchant_service.RestaurantAuth())

			///////////////商店
			merchant_session.GET("/restaurants",
				v1.GetRestaurants)
			merchant_session.POST(
				"restaurant",
				v1.CreateRestaurant,
			)
			merchant_restaurant.DELETE(
				"",
				v1.DeleteRestaurant,
			)
			merchant_restaurant.PUT(
				"",
				v1.UpdateRestaurant,
			)

			///////////////商店状态
			merchant_restaurant.PUT(
				"/status/:status",
				v1.SetRestaurantStatus,
			)
			merchant_restaurant.GET(
				"/status",
				v1.GetRestaurantStatus,
			)

			///////////////菜品类别
			merchant_restaurant.GET(
				"/categories",
				v1.GetCategories,
			)
			merchant_restaurant.POST(
				"/category",
				v1.CreateCategory,
			)
			merchant_session.PUT(
				"/category/:category_id",
				// 这里需要验证更新的category是否是商家自己的店铺的
				merchant_service.CategoryAuth(),
				v1.UpdateCategory,
			)
			merchant_session.DELETE(
				"/category/:category_id",
				merchant_service.CategoryAuth(),
				v1.DeleteCategory,
			)

			///////////////口味
			merchant_restaurant.POST(
				"/flavor/:name",
				v1.CreateFlavor,
			)
			merchant_session.DELETE(
				"/flavor/:flavor_id",
				v1.DeleteFlavor,
			)
			merchant_session.PUT(
				"/flavor/:flavor_id/name/:name",
				v1.UpdateFlavor,
			)
			merchant_restaurant.GET(
				"/flavors",
				v1.GetRestaurantFlavors,
			)

			///////////////菜品的口味
			// form中发送需要加入的
			merchant_session.POST(
				"/dish/:dish_id/flavors/add",
				// merchant_service.DishAuth(),
				v1.AddDishFlavor,
			)
			merchant_session.POST(
				"/dish/:dish_id/flavors/delete",
				// merchant_service.DishAuth(),
				v1.DeleteDishFlavor,
			)
			merchant_session.GET(
				"/dish/:dish_id/flavors",
				// merchant_service.DishAuth(),
				v1.GetDishFlavor,
			)

			///////////////菜品
			// NOTE:这里需要验证更新的dish是否属于商家
			merchant_session.DELETE(
				"/dish/:dish_id",
				// merchant_service.DishAuth(),
				v1.DeleteDish,
			)
			merchant_session.PUT(
				"/dish/:dish_id",
				// merchant_service.DishAuth(),
				v1.UpdateDish,
			)
			merchant_session.POST(
				"/category/:category_id/dishes/add",
				v1.AddCategoryDish,
			)
			merchant_session.POST(
				"/category/:category_id/dishes/delete",
				v1.DeleteCategoryDish,
			)
			// NOTE:菜品直接归属在restaurant名下
			merchant_restaurant.POST(
				"/dish",
				v1.CreateDish,
			)
			merchant_restaurant.GET(
				"/dish",
				v1.GetDishes,
			)

			// NOTE: license的图片静态文件路由, 对餐品图片的访问不进行鉴权
			dish_image_path := setting.AppSetting.DishImageStorePath
			log.Infof("Serving Static File: [%s]", dish_image_path)
			merchant_session.Static("/dish/image", dish_image_path)
		}

	}

	//////////////////////////////////////////////////////
	//////////////////////////////////////////////////////
	//////////////////////////////////////////////////////
	//////////////////////////////////////////////////////

	{
		// middleware redis session
		admin_store_v1, err := redis.NewStore(
			setting.RedisSetting.MaxIdle,
			"tcp",
			setting.RedisSetting.Host,
			setting.RedisSetting.Password,
			[]byte(setting.RedisSetting.Secret),
		)
		if err != nil {
			log.Fatal(err)
		}
		admin_store_v1.Options(sessions.Options{
			Path:     "/api/v1/admin",
			MaxAge:   setting.AppSetting.AdminAliveMinute * 60,
			Secure:   true, // 仅通过HTTPS传输Cookie
			HttpOnly: true, // 禁止通过JavaScript访问Cookie
			SameSite: http.SameSiteDefaultMode,
		})
		admin_session_handler := sessions.Sessions("AdminSession", admin_store_v1)

		admin_account_v1 := apiv1.Group("/admin")
		admin_account_v1.Use(admin_session_handler)
		admin_account_v1.POST("/create", api.AdminCreate)
		admin_account_v1.DELETE("/delete", api.AdminDelete)
		admin_account_v1.POST("/login", api.AdminLogin)
		admin_account_v1.POST("/logout", auth.AdminAuth(), api.AdminLogout)

		{
			admin_session_v1 := apiv1.Group("/admin")
			admin_session_v1.Use(admin_session_handler)
			admin_session_v1.Use(auth.AdminAuth())
			admin_session_v1.GET("/login-status", api.AdminLoginStatus)
			admin_session_v1.PUT("/change-password", api.AdminChangePassword)
			admin_session_v1.POST("/merchant/create",
				v1.CreateMerchant)
			admin_session_v1.POST("/merchant/delete",
				v1.DeleteMerchant)
			admin_session_v1.POST("/merchant/:merchant_id/disable",
				v1.DisableMerchant)
			admin_session_v1.POST("/merchant/:merchant_id/enable",
				v1.EnableMerchant)
			admin_session_v1.GET("/merchant-application/:page",
				v1.GetMerchantApplication)
			admin_session_v1.GET("/merchants/:page",
				v1.GetMerchants)
			admin_session_v1.PUT("/merchant-application/:application_id/approve",
				v1.ApproveMerchantApplication)
			admin_session_v1.PUT("/merchant-application/:application_id/disapprove",
				v1.DisapproveMerchantApplication)

			// NOTE: license的图片静态文件路由
			license_path := setting.AppSetting.LicenseStorePath
			log.Infof("Serving Static File: [%s]", license_path)
			admin_session_v1.Static("/merchant-application/license", license_path)
		}
	}

	return r
}
