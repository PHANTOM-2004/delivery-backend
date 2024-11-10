package docs

import (
	"bytes"

	"github.com/go-openapi/runtime"
)

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
// swagger:route PUT /api/v1/merchant/jwt/dish/{dish_id} v1-merchant merchant_update_dish
// 更新某一个菜品
// PS:分类中包含菜品项，所以实际上获得菜品的接口也是这个
// responses:
// 200: COMMON

// =============================================================
// swagger:route POST /api/v1/merchant/jwt/category/{category_id}/dish v1-merchant merchant_create_dish
// 创建某一个分类下的菜品
// responses:
// 200: COMMON

//swagger:parameters merchant_delete_dish
type DeleteDish struct {
	//in:path
	DishID uint `json:"dish_id"`
}

// =============================================================
// swagger:route DELETE /api/v1/merchant/jwt/dish/{dish_id} v1-merchant merchant_delete_dish
// 删除某一个菜品
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
// swagger:route GET /api/v1/merchant/jwt/dish/image/{*filepath} v1-merchant merchant_get_dish_image
// 请求得到菜品照片; 如果路径内容不存在则对应httpcode=404
// responses:
// 200: dish_image
