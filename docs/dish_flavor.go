package docs

import "delivery-backend/models"

//swagger:parameters merchant_add_dish_flavor merchant_delete_dish_flavor
type FlavorAssDishIDRequest struct {
	//in:path
	DishID uint `json:"dish_id"`
	// 需要添加的flavor id数组
	//in:formData
	//required:true
	Flavors []uint `json:"flavors"`
}

// =============================================================
// swagger:route POST  /api/v1/merchant/dish/{dish_id}/flavors/add v1-merchant merchant_add_dish_flavor
// 为某一个菜品添加口味
// 支持一次传递多个口味
// responses:
// 200: COMMON

// =============================================================
// swagger:route POST  /api/v1/merchant/dish/{dish_id}/flavors/delete v1-merchant merchant_delete_dish_flavor
// 为一个菜品删除口味
// 支持一次传递多个口味
// responses:
// 200: COMMON

//swagger:response merchant_get_dish_flavor_response
type DishFlavorResponse struct {
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

//swagger:parameters merchant_get_dish_flavor
type GetDishFlavorRequest struct {
	//in:path
	DishID uint `json:"dish_id"`
}

// =============================================================
// swagger:route GET  /api/v1/merchant/dish/{dish_id}/flavors v1-merchant merchant_get_dish_flavor
// 获得一个dish的flavor
// responses:
// 200: merchant_get_dish_flavor_response
