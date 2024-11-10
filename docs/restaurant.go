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
// swagger:route GET /api/v1/merchant/jwt/restuarants v1-merchant merchant_get_restaurants
// 商家GET名下的所有店铺
// responses:
// 200: merchant_get_restaurants

// swagger:parameters merchant_get_restaurant_status merchant_get_restaurant_flavors
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
	Status uint8 `json:"status"`

	// required: true
	// in: path
	RestaurantID uint `json:"restaurant_id"`
}

// =============================================================
// swagger:route PUT /api/v1/merchant/jwt/restuarant/{restaurant_id}/status/{status} v1-merchant merchant_set_restaurant_status
//
// 设置商家某个店铺的状态
// 设置0, 代表商家手动关闭店铺；设置1, 代表店铺开启。
// PS:目前只支持商家手动设置店铺状态。
// responses:
// 200: COMMON

// swagger:response merchant_get_restaurant_flavors_response
type RestaurnatDishesResponse struct {
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
// swagger:route GET /api/v1/merchant/jwt/restaurant/{restaurant_id}/flavors v1-merchant merchant_get_restaurant_flavors
// 获得一个餐馆的所有口味
// responses:
// 200: merchant_get_restaurant_flavors_response
