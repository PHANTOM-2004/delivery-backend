package docs

//swagger:parameters merchant_add_category_dish merchant_delete_category_dish
type DishAssCategoryRequest struct {
	//in:path
	CategoryID uint `json:"category_id"`
	// 需要添加的dish id数组
	//in:formData
	//required:true
	Dishes []uint `json:"dishes"`
}

// =============================================================
// swagger:route POST  /api/v1/merchant/category/{category_id}/dishes/add v1-merchant merchant_add_category_dish
// 为某一个菜品添加口味
// 支持一次传递多个口味
// responses:
// 200: COMMON

// =============================================================
// swagger:route POST  /api/v1/merchant/category/{category_id}/dishes/delete v1-merchant merchant_delete_category_dish
// 为一个菜品删除口味
// 支持一次传递多个口味
// responses:
// 200: COMMON
