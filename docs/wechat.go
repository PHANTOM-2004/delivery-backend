package docs

import (
	"delivery-backend/models"
	v1 "delivery-backend/routers/api/v1"
)

// swagger:parameters wechat_login
type WXLoginRequest struct {
	//required:true
	//in:body
	Object struct {
		Code string `json:"code"`
	} `json:"object"`
}

// swagger:response wechat_login
type WXLoginResponse struct {
	// in:body
	Body struct {
		// Required:true
		// Example: 10000
		ECode int `json:"ecode"`
		// Example: 管理员不存在
		// error message
		// Required:true
		Msg string `json:"msg"`
		// Required:true
		// data to get
		// Required:true
		Body struct {
			// required: true
			SessionID string `json:"session_id"`
			// 当该用户不是新用户的时候返回
			Info v1.UserInfoRequest `json:"info"`
		} `json:"data"`
	}
}

// =============================================================
// swagger:route POST /api/v1/wx/login v1-wechat wechat_login
// 微信登录，会返回sessionid
// 注意传递数据是json格式。注意info字段的返回条件。
// responses:
// 200: wechat_login
//
//

// swagger:parameters customer_get_cdf
type CustomerGetCDFRequest struct {
	// in:path
	// Required: true
	RestaurantID uint `json:"restaurant_id"`
}

// swagger:response customer_get_cdf
type CustomerGetCDFResponse struct {
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
			Category []models.Category `json:"categories"`
		} `json:"data"`
	}
}

// =============================================================
// swagger:route GET /api/v1/wx/customer/restaurant/{restaurant_id}/categories/dishes v1-wechat customer_get_cdf
// 返回商家某个商店的所有分类,以及分类下的所有菜品，以及菜品下的所有口味
// responses:
// 200: customer_get_cdf
