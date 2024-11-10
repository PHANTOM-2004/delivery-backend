package docs

import "delivery-backend/models"

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

// swagger:parameters merchant_get_categories merchant_get_category merchant_create_category merchant_delete_category
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

// swagger:parameters merchant_get_category merchant_update_category merchant_create_dish
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
// swagger:route DELETE /api/v1/merchant/jwt/category/{category_id} v1-merchant merchant_delete_category
// 删除某个分类, 注意删除的不是自己的情况
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
// swagger:route PUT /api/v1/merchant/jwt/category/{category_id} v1-merchant merchant_update_category
// 更新某一个菜品分类分类
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
// swagger:route POST /api/v1/merchant/jwt/restuarant/{restaurant_id}/category v1-merchant merchant_create_category
// 创建一个餐馆下的菜品分类
// responses:
// 200: COMMON
