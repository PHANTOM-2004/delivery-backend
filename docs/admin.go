package docs

import "github.com/go-openapi/runtime"

// =============================================================
// 用于登入与认证的参数，注意后端处理的最大范围, 否则会返回错误
// swagger:parameters admin_login admin_create admin_delete
type AdminAccountParam struct {
	// 账户名：最小长度10，最大长度30
	// required: true
	// in: formData
	Account string `json:"account"`
}

// swagger:parameters admin_login admin_create admin_change_password
type AdminPasswordParam struct {
	// 密码：最小长度15, 最大长度30
	// required: true
	// in: formData
	Password string `json:"password"`
}

// SuperToken, 用于管理员敏感操作
// swagger:parameters  admin_create admin_delete
type SuperToken struct {
	// 超级管理员的token, 该token源于定期生成或者已持有密钥(如见app.ini)
	// 当正确的token为空时，该接口被禁用
	// required: true
	SuperToken string `json:"super_token"`
}

// swagger:parameters  admin_create
type AdminName struct {
	// 管理员姓名, 最小长度2, 最大长度20
	// required: true
	// in: formData
	AdminName string `json:"admin_name"`
}

//=============================================================
// swagger:route POST /api/v1/admin/login v1-admin admin_login
// 登入的身份认证
// (*) 在cookie中设置session_id
// PS: 通过postform发送参数, 否则会认证错误.
// responses:
// 200: COMMON
//

//=============================================================
// swagger:route GET /api/v1/admin/login-status v1-admin admin_login_status
// 请求管理员的登陆状态
// (1) 如果登陆，同时会返回已登录的account，data字段中有一个key为account
// (2) 在判断httpcode的基础上(httpcode != 401)，只需要判断业务逻辑码是否是SUCCESS即可,不存在error时意味着处于登录状态
// (3) 可能存在角色错误或者未登陆的错误码
// responses:
// 200: COMMON

//=============================================================
// swagger:route POST /api/v1/admin/create v1-admin admin_create
// 创建管理员账户，该api只允许测试时以及部署时运维调用
// PS: 通过url发送参数
// responses:
// 200: COMMON

// =============================================================
// swagger:route DELETE /api/v1/admin/delete v1-admin admin_delete
// 删除管理员账户，该api只允许测试时以及部署时运维调用
// PS: 通过url发送参数，删除账户不存在时也会删除成功，但是会返回信息提示不存在
// responses:
// 200: COMMON

// =============================================================
// swagger:route POST /api/v1/admin/logout v1-admin admin_logout
// 注销管理员账户
// 在非法请求发出时（管理员不处于登入状态）会返回错误信息。
// responses:
// 200: COMMON

// =============================================================
// swagger:route PUT /api/v1/admin/change-password v1-admin admin_change_password
// 管理员修改密码
// 在非法请求发出时（管理员不处于登入状态）会返回错误信息。
// PS: 通过postform形式传递密码, 不要使用url传参。
// responses:
// 200: COMMON

// swagger:parameters admin_approve_application admin_disapprove_application
type MerchantApplicationID struct {
	// 申请表id
	// in: path
	ID string `json:"application_id"`
}

// =============================================================
// swagger:route PUT /api/v1/admin/merchant-application/{application_id}/approve v1-admin admin_approve_application
// 管理员通过商家发起的申请
// 如果同意，会判断是否存在关联的商家账号，如果存在必然是封禁状态，那么会对状态进行解禁; 如果不存在关联的商家账号，那么会为该申请创建一个新账号。
// PS: 目前新账号会放在日志中，后续新账号会发送到申请邮箱中。
// responses:
// 200: COMMON

// =============================================================
// swagger:route PUT /api/v1/admin/merchant-application/{application_id}/disapprove v1-admin admin_disapprove_application
// 管理员拒绝商家发起的申请
// 如果拒绝, 仍然会判断是否有关联的商家账号，如果没有关联的商家账号，那么仅拒绝申请。
// 如果有关联的商家账号，还会把这个关联的商家账号进行冻结。
// 在非法请求发出时（管理员不处于登入状态）会返回错误信息。
// responses:
// 200: COMMON

// swagger:parameters  admin_get_merchant_license
type MerchantLicenseRequest struct {
	// example: merchant-license-02fd3ce1-fdcb-4a30-94a4-db3f9c241871.png
	// in: path
	ID string `json:"*filepath"`
}

// swagger:response merchant_license
type MerchantLicenseResponse struct {
	// 执照图片
	// Required:true
	License runtime.File
}

// =============================================================
// swagger:route GET /api/v1/admin/merchant-application/license/{*file_path} v1-admin admin_get_merchant_license
// 请求得到商家执照; 如果路径内容不存在则对应httpcode=404
// responses:
// 200: merchant_license
