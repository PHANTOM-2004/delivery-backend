package docs

import (
	"delivery-backend/models"
)

//swagger:parameters customer_create_order
type CreateOrderRequest struct {
	//in:path
	RestaurantID uint `json:"restaurant_id"`
	// in: body
	Body struct {
		// 最大长度100
		// required:true
		Address string `json:"address"`
		// 最大长度20
		// required:true
		CustomerName string `json:"customer"`
		// e164规范
		// required:true
		PhoneNumber string `json:"phone_number"`
	}
}

// =============================================================
// swagger:route POST /api/v1/wx/customer/order/restaurant/{restaurant_id} v1-wechat customer_create_order
// 创建订单
// responses:
// 200: COMMON

//swagger:response customer_get_orders
type GetUserOrderResponse struct {
	//in:body
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
			Cart []models.Order `json:"orders"`
		} `json:"data"`
	}
}

// =============================================================
// swagger:route GET /api/v1/wx/customer/orders v1-wechat customer_create_order
// 创建订单
// responses:
// 200: customer_get_orders
