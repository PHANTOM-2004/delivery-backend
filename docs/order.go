package docs

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
