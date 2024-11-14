package docs

import "delivery-backend/models"

// =============================================================
// 用于登入与认证的参数，注意后端处理的最大范围, 否则会返回错误
// swagger:parameters merchant_login merchant_create merchant_delete
type MerchantAccountParam struct {
	// 账户名：最小长度6，最大长度30
	// in: formData
	// required: true
	Account string `json:"account"`
}

// swagger:parameters merchant_login merchant_create merchant_change_password
type MerchantPasswordParam struct {
	// 密码：最小长度8, 最大长度30
	// in: formData
	// required: true
	Password string `json:"password"`
}

// swagger:parameters  merchant_create
type MerchantName struct {
	// 商家姓名, 最小长度2, 最大长度20
	// in: formData
	// required: true
	MerchantName string `json:"merchant_name"`
}

// swagger:parameters merchant_create
type MerchantPhoneNumber struct {
	// 商家手机号, 遵照E.164规范
	// in: formData
	// required: true
	PhoneNumber string `json:"phone_number"`
}

// swagger:parameters merchant_create
type MerchantCreateApplicationID struct {
	// 商家对应的application_id
	// in: formData
	// required: true
	ID string `json:"merchant_application_id"`
}

//=============================================================
// swagger:route POST /api/v1/merchant/login v1-merchant merchant_login
// 登入的身份认证
// (*) 在cookie中设置session_id
// PS: 通过postform发送参数, 否则会认证错误.
// responses:
// 200: COMMON
//

//=============================================================
// swagger:route GET /api/v1/merchant/login-status v1-merchant merchant_login_status
// 请求商家的登陆状态
// (1) 如果登陆，同时会返回已登录的account，data字段中有一个key为account
// (2) 在判断httpcode的基础上(httpcode != 401)，只需要判断业务逻辑码是否是SUCCESS即可,不存在error时意味着处于登录状态
// (3) 可能存在角色错误或者未登陆的错误码
// responses:
// 200: COMMON

//=============================================================
// swagger:route POST /api/v1/admin/merchant/create v1-admin merchant_create
// 创建商家账户
// 该api将会由管理员调用，接受商家的账号申请创建
// PS: 通过postform发送参数
// responses:
// 200: COMMON

// =============================================================
// swagger:route POST /api/v1/admin/merchant/delete v1-admin merchant_delete
// 删除商家账户
// 该api将由管理员调用，删除商家；或者商家自己注销账号。
// PS: 通过postform发送参数，删除账户不存在时也会删除成功，但是会返回信息提示不存在
// responses:
// 200: COMMON

// =============================================================
// swagger:route POST /api/v1/merchant/logout v1-merchant merchant_logout
// 注销商家账户
// 在非法请求发出时（商家不处于登入状态）会返回错误信息。
// responses:
// 200: COMMON

// =============================================================
// swagger:route PUT /api/v1/merchant/change-password v1-merchant merchant_change_password
// 商家修改密码
// 在非法请求发出时（商家不处于登入状态）会返回错误信息。
// PS: 通过postform形式传递密码, 不要使用url传参。
// responses:
// 200:COMMON

// swagger:response merchant_information
type MerchantInfoResponse struct {
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
		// In:body
		Data struct {
			Merchant models.Merchant `json:"merchant"`
		} `json:"data"`
	}
}

// =============================================================
// swagger:route GET /api/v1/merchant/info v1-merchant merchant_information
// 商家获得个人信息
// responses:
// 200:merchant_information
