package docs

import (
	"bytes"
	"delivery-backend/models"

	"github.com/go-openapi/runtime"
)

// =============================================================
// 用于登入与认证的参数，注意后端处理的最大范围, 否则会返回错误
// swagger:parameters merchant_login merchant_create merchant_delete
type MerchantAccountParam struct {
	// 账户名：最小长度6，最大长度30
	// in: formData
	// required: true
	Account string `json:"account"`
}

// swagger:parameters merchant_login merchant_create merchant_change_password
type MerchantPasswordParam struct {
	// 密码：最小长度8, 最大长度30
	// in: formData
	// required: true
	Password string `json:"password"`
}

// swagger:parameters  merchant_create
type MerchantName struct {
	// 商家姓名, 最小长度2, 最大长度20
	// in: formData
	// required: true
	MerchantName string `json:"merchant_name"`
}

// swagger:parameters merchant_create
type MerchantPhoneNumber struct {
	// 商家手机号, 遵照E.164规范
	// in: formData
	// required: true
	PhoneNumber string `json:"phone_number"`
}

// swagger:parameters merchant_create
type MerchantCreateApplicationID struct {
	// 商家对应的application_id
	// in: formData
	// required: true
	ID string `json:"merchant_application_id"`
}

//=============================================================
// swagger:route GET /api/v1/merchant/jwt/auth v1-merchant merchant_auth
// 请求获得access_token
// (1) 通过refresh_token获取access_token
// (2) 注意错误码，如果出现refresh_token过期说明需要重新登录
// responses:
// 200: COMMON

//=============================================================
// swagger:route POST /merchant/login merchant merchant_login
// 登入的身份认证
// (1) 返回access_token, refresh_token
// (2) 其中access_token是短期的有效token, refresh_token是长期有效token, 后者用于刷新acess_token
// PS: 通过postform发送参数, 否则会认证错误.
// responses:
// 200: COMMON
//

//=============================================================
// swagger:route GET /api/v1/merchant/jwt/login-status v1-merchant merchant_login_status
// 请求商家的登陆状态
// (1) 如果登陆，同时会返回已登录的account，data字段中有一个key为account
// (2) 在判断httpcode的基础上(httpcode != 401)，只需要判断业务逻辑码是否是SUCCESS即可,不存在error时意味着处于登录状态
// (3) 可能存在角色错误或者未登陆的错误码
// responses:
// 200: COMMON

//=============================================================
// swagger:route POST /api/v1/admin/jwt/merchant-create v1-admin merchant_create
// 创建商家账户
// 该api将会由管理员调用，接受商家的账号申请创建
// PS: 通过postform发送参数
// responses:
// 200: COMMON

// =============================================================
// swagger:route DELETE /api/v1/admin/jwt/merchant-delete v1-admin merchant_delete
// 删除商家账户
// 该api将由管理员调用，删除商家；或者商家自己注销账号。
// PS: 通过postform发送参数，删除账户不存在时也会删除成功，但是会返回信息提示不存在
// responses:
// 200: COMMON

// =============================================================
// swagger:route POST /merchant/logout merchant merchant_logout
// 注销商家账户
// 在非法请求发出时（商家不处于登入状态）会返回错误信息。
// responses:
// 200: COMMON

// =============================================================
// swagger:route PUT /api/v1/merchant/jwt/change-password v1-merchant merchant_change_password
// 商家修改密码
// 在非法请求发出时（商家不处于登入状态）会返回错误信息。
// PS: 通过postform形式传递密码, 不要使用url传参。
// responses:
// 200:COMMON

// swagger:response merchant_get_restaurants
type MerchantRestaurantsResponse struct {
	// in:body
	Body struct {
		// Required:true
		// Example: 200
		ECode int `json:"ecode"`
		// Example: ok
		// error message
		// Required:true
		Msg string `json:"msg"`
		// Required:true
		// data to get
		// Required:true
		Data struct {
			// in:body
			Restaurants []models.Restaurant `json:"restaurants"`
		} `json:"data"`
	}
}

// =============================================================
// swagger:route GET /api/v1/merchant/jwt/restuarants v1-merchant merchant_get_restaurants
// 商家GET名下的所有店铺
// responses:
// 200: merchant_get_restaurants

// swagger:parameters merchant_get_restaurant_status
type MerchantRestaurantStatusGetRequest struct {
	// required: true
	//in:path
	RestaurantID uint `json:"restaurant_id"`
}

// =============================================================
// swagger:route GET /api/v1/merchant/jwt/restuarant/{restaurant_id}/status v1-merchant merchant_get_restaurant_status
//
// 获取商家某个店铺的状态
// 如果返回0, 代表商家手动关闭店铺；如果返回1, 代表店铺开启。
// PS:目前只支持商家手动设置店铺状态。
// responses:
// 200: COMMON

// swagger:parameters merchant_set_restaurant_status
type MerchantRestaurantStatusSetRequest struct {
	// required: true
	// in: path
	Status uint `json:"status"`

	// required: true
	// in: path
	RestaurantID uint8 `json:"restaurant_id"`
}

// =============================================================
// swagger:route PUT /api/v1/merchant/jwt/restuarant/{restaurant_id}/status/{status} v1-merchant merchant_set_restaurant_status
//
// 设置商家某个店铺的状态
// 设置0, 代表商家手动关闭店铺；设置1, 代表店铺开启。
// PS:目前只支持商家手动设置店铺状态。
// responses:
// 200: COMMON

// swagger:response merchant_get_categories_response
type MerchantGetCategoriesResponse struct {
	// in:body
	Body struct {
		// Required:true
		// Example: 200
		ECode int `json:"ecode"`
		// Example: ok
		// error message
		// Required:true
		Msg string `json:"msg"`
		// Required:true
		// data to get
		// Required:true
		Data struct {
			// in:body
			Categories []models.Category `json:"categories"`
		} `json:"data"`
	}
}

// swagger:parameters merchant_get_categories merchant_get_category merchant_update_category merchant_create_category merchant_update_dish merchant_create_dish
type MerchantRestaurantRequest struct {
	// in:path
	// Required: true
	RestaurantID uint `json:"restaurant_id"`
}

// =============================================================
// swagger:route GET /api/v1/merchant/jwt/restuarant/{restaurant_id}/categories v1-merchant merchant_get_categories
// 返回商家某个商店的所有分类
// PS:分类中包含菜品项，所以实际上获得菜品的接口也是这个
// responses:
// 200: merchant_get_categories_response

// swagger:parameters merchant_get_category merchant_update_category merchant_update_dish merchant_create_dish
type MerchantGetCategoryRequest struct {
	// in:path
	// Required: true
	CategoryID uint `json:"category_id"`
}

// swagger:response merchant_get_category_response
type MerchantGetCategoryResponse struct {
	// in:body
	Body struct {
		// Required:true
		// Example: 200
		ECode int `json:"ecode"`
		// Example: ok
		// error message
		// Required:true
		Msg string `json:"msg"`
		// Required:true
		// data to get
		// Required:true
		Data struct {
			// in:body
			Category models.Category `json:"categories"`
		} `json:"data"`
	}
}

// =============================================================
// swagger:route GET /api/v1/merchant/jwt/restuarant/{restaurant_id}/category/{category_id} v1-merchant merchant_get_category
// 返回商家某个商店的所有分类
// PS:分类中包含菜品项，所以实际上获得菜品的接口也是这个
// responses:
// 200: merchant_get_category_response

// swagger:parameters merchant_update_category
type MerchantUpdateCategoryRequest struct {
	//in:path
	//required:true
	ID uint `json:"category_id"`
	//最大长度: 30bytes
	//in:formData
	Name string `json:"name"`
	//0代表菜品，1代表套餐
	//in:formData
	Type uint8 `json:"type"`
	//in:formData
	Sort uint16 `json:"sort"`
}

// =============================================================
// swagger:route PUT /api/v1/merchant/jwt/restuarant/{restaurant_id}/category/{category_id}/update v1-merchant merchant_update_category
// 返回商家某个商店的所有分类
// PS:分类中包含菜品项，所以实际上获得菜品的接口也是这个
// responses:
// 200: COMMON

// swagger:parameters merchant_create_category
type MerchantCreateCategoryRequest struct {
	//最大长度: 30bytes
	//in:formData
	//required:true
	Name string `json:"name"`
	//0代表菜品，1代表套餐
	//in:formData
	//required:true
	Type uint8 `json:"type"`
	//in:formData
	//required:true
	Sort uint16 `json:"sort"`
}

// =============================================================
// swagger:route POST /api/v1/merchant/jwt/restuarant/{restaurant_id}/category/{category_id}/create v1-merchant merchant_create_category
// 返回商家某个商店的所有分类
// PS:分类中包含菜品项，所以实际上获得菜品的接口也是这个
// responses:
// 200: COMMON

// swagger:parameters merchant_update_dish
type MerchantUpdateDishRequest struct {
	//in:path
	//required:true
	ID uint `json:"dish_id"`
	//最大长度: 30bytes
	//in:formData
	Name string `json:"name"`
	//使用整数存储价格，精确到分
	//in:formData
	Price uint `json:"price"`
	//in:formData
	Sort uint16 `json:"sort"`
	//最大长度: 50 bytes
	//in:formData
	Description string `json:"description"`
	// 文件存在最大大小，参考配置文件
	// multipart/form data, 上传一个文件; 接受格式.png,.jpeg,.jpg
	// in: formData
	// swagger:file
	Image *bytes.Buffer `json:"image"`
}

// =============================================================
// swagger:route PUT /api/v1/merchant/jwt/restuarant/{restaurant_id}/category/{category_id}/dish/{dish_id}/update v1-merchant merchant_update_dish
// 返回商家某个商店的所有分类
// PS:分类中包含菜品项，所以实际上获得菜品的接口也是这个
// responses:
// 200: COMMON

// swagger:parameters merchant_create_dish
type MerchantCreateDishRequest struct {
	//in:path
	//required:true
	ID uint `json:"dish_id"`
	//最大长度: 30bytes
	//in:formData
	// requires:true
	Name string `json:"name"`
	//使用整数存储价格，精确到分
	// requires:true
	//in:formData
	Price uint `json:"price"`
	//默认值是0
	//in:formData
	Sort uint16 `json:"sort"`
	//最大长度: 50 bytes
	// requires:true
	//in:formData
	Description string `json:"description"`
	// 文件存在最大大小，参考配置文件
	// multipart/form data, 上传一个文件; 接受格式.png,.jpeg,.jpg
	// in: formData
	// requires:true
	// swagger:file
	Image *bytes.Buffer `json:"image"`
}

// =============================================================
// swagger:route POST /api/v1/merchant/jwt/restuarant/{restaurant_id}/category/{category_id}/dish/create v1-merchant merchant_create_dish
// 返回商家某个商店的所有分类
// PS:分类中包含菜品项，所以实际上获得菜品的接口也是这个
// responses:
// 200: COMMON

// swagger:parameters merchant_get_dish_image
type DishImageRequest struct {
	// example: dish-02fd3ce1-fdcb-4a30-94a4-db3f9c241871.png
	// in: path
	ID string `json:"*filepath"`
}

// swagger:response dish_image
type DishImageResponse struct {
	// 执照图片
	// Required:true
	Image runtime.File
}

// =============================================================
// swagger:route GET /api/v1/merchant/jwt/dish/{*filepath} v1-merchant merchant_get_dish_image
// 请求得到菜品照片; 如果路径内容不存在则对应httpcode=404
// responses:
// 200: dish_image
