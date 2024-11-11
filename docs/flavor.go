package docs

//swagger:parameters merchant_create_flavor
type RestaurantFlavorsRequest struct {
	// in: path
	RestaurantID uint `json:"restaurant_id"`
}

//swagger:parameters merchant_create_flavor merchant_update_flavor
type UpdateFlavorName struct {
	// in: path
	Name string `json:"name"`
}

//swagger:parameters merchant_update_flavor merchant_delete_flavor
type FlavorIDRequest struct {
	//in:path
	FlavorID uint `json:"flavor_id"`
}

// =============================================================
// swagger:route POST /api/v1/merchant/restaurant/{restaurant_id}/flavor/{name} v1-merchant merchant_create_flavor
// 创建某一个口味
// responses:
// 200: COMMON

// =============================================================
// swagger:route DELETE /api/v1/merchant/flavor/{flavor_id} v1-merchant merchant_delete_flavor
// 删除某一个口味
// responses:
// 200: COMMON

// =============================================================
// swagger:route PUT /api/v1/merchant/flavor/{flavor_id}/name/{name} v1-merchant merchant_update_flavor
// 更新某一个口味
// responses:
// 200: COMMON
