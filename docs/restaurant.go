package docs

import "delivery-backend/models"

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
// swagger:route GET /api/v1/merchant/restaurants v1-merchant merchant_get_restaurants
// 商家GET名下的所有店铺
// responses:
// 200: merchant_get_restaurants

// swagger:parameters merchant_get_restaurant_status merchant_get_restaurant_flavors merchant_delete_restaurant merchant_update_restaurant
type MerchantRestaurantStatusGetRequest struct {
	// required: true
	//in:path
	RestaurantID uint `json:"restaurant_id"`
}

// =============================================================
// swagger:route GET /api/v1/merchant/restaurant/{restaurant_id}/status v1-merchant merchant_get_restaurant_status
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
	Status uint8 `json:"status"`

	// required: true
	// in: path
	RestaurantID uint `json:"restaurant_id"`
}

// =============================================================
// swagger:route PUT /api/v1/merchant/restaurant/{restaurant_id}/status/{status} v1-merchant merchant_set_restaurant_status
//
// 设置商家某个店铺的状态
// 设置0, 代表商家手动关闭店铺；设置1, 代表店铺开启。
// PS:目前只支持商家手动设置店铺状态。
// responses:
// 200: COMMON

// swagger:response merchant_get_restaurant_flavors_response
type RestaurantDishesResponse struct {
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
			Flavors []models.Flavor `json:"flavors"`
		} `json:"data"`
	}
}

// =============================================================
// swagger:route GET /api/v1/merchant/restaurant/{restaurant_id}/flavors v1-merchant merchant_get_restaurant_flavors
// 获得一个餐馆的所有口味
// responses:
// 200: merchant_get_restaurant_flavors_response

//swagger:parameters merchant_create_restaurant
type RestaurantCreateRequest struct {
	// in: formData
	// Required: true
	RestaurantName string `json:"restaurant_name"`
	// 店铺的地址
	// in: formData
	// Required: true
	Address string `json:"address"`
	// 商铺简介
	// in: formData
	// Required: true
	Description string `json:"description"`
	// 最小起送金额,使用整数存储,默认存储到分
	// in: formData
	// Required: true
	MinimumDeliveryAmount uint `json:"minimum_delivery_amount"`
}

// =============================================================
// swagger:route POST /api/v1/merchant/restaurant v1-merchant merchant_create_restaurant
// 创建一个餐厅
// responses:
// 200: COMMON

//swagger:parameters merchant_update_restaurant
type RestaurantUpdateRequest struct {
	// in: formData
	RestaurantName string `json:"restaurant_name"`
	// 店铺的地址
	// in: formData
	Address string `json:"address"`
	// 商铺简介
	// in: formData
	Description string `json:"description"`
	// 最小起送金额,使用整数存储,默认存储到分
	// in: formData
	MinimumDeliveryAmount uint `json:"minimum_delivery_amount"`
}

// =============================================================
// swagger:route PUT /api/v1/merchant/restaurant/{restaurant_id} v1-merchant merchant_update_restaurant
// 更新一个餐厅
// responses:
// 200: COMMON

// =============================================================
// swagger:route DELETE /api/v1/merchant/restaurant/{restaurant_id} v1-merchant merchant_delete_restaurant
// 删除一个餐厅
// responses:
// 200: COMMON
