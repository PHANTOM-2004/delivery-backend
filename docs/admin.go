package docs

// =============================================================
// 用于登入与认证的参数，注意后端处理的最大范围, 否则会返回错误
// swagger:parameters admin_login admin_create admin_delete
type AdminAccountParam struct {
	// 账户名：最小长度10，最大长度30
	// required: true
	Account string `json:"account"`
}

// swagger:parameters admin_login admin_create admin_change_password
type AdminPasswordParam struct {
	// 密码：最小长度15, 最大长度30
	// required: true
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

// [UNUSED swagger]:parameters  admin_change_password
type AccessToken struct {
	// 部分api调用的token
	// required: true
	AccessToken string `json:"access_token"`
}

// swagger:parameters  admin_create
type AdminName struct {
	// 管理员姓名, 最小长度2, 最大长度20
	// required: true
	AdminName string `json:"admin_name"`
}

//=============================================================
// swagger:route GET /api/v1/admin/jwt/auth v1-admin-jwt admin_auth
// 请求获得access_token
// (1) 通过refresh_token获取access_token
// (2) 注意错误码，如果出现refresh_token过期说明需要重新登录
// responses:
// 200: COMMON

//=============================================================
// swagger:route POST /admin/login admin admin_login
// 登入的身份认证
// (1) 返回access_token, refresh_token
// (2) 其中access_token是短期的有效token, refresh_token是长期有效token, 后者用于刷新acess_token
// PS: 通过postform发送参数, 否则会认证错误.
// responses:
// 200: COMMON
//

//=============================================================
// swagger:route GET /api/v1/admin/jwt/login-status v1-admin-jwt admin_login_status
// 请求管理员的登陆状态
// (1) 如果登陆，同时会返回已登录的account，data字段中有一个key为account
// (2) 在判断httpcode的基础上(httpcode != 401)，只需要判断业务逻辑码是否是SUCCESS即可,不存在error时意味着处于登录状态
// (3) 可能存在角色错误或者未登陆的错误码
// responses:
// 200: COMMON

//=============================================================
// swagger:route POST /admin/create admin admin_create
// 创建管理员账户，该api只允许测试时以及部署时运维调用
// PS: 通过url发送参数
// responses:
// 200: COMMON

// =============================================================
// swagger:route DELETE /admin/delete admin admin_delete
// 删除管理员账户，该api只允许测试时以及部署时运维调用
// PS: 通过url发送参数，删除账户不存在时也会删除成功，但是会返回信息提示不存在
// responses:
// 200: COMMON

// =============================================================
// swagger:route POST /admin/logout admin admin_logout
// 注销管理员账户
// 在非法请求发出时（管理员不处于登入状态）会返回错误信息。
// responses:
// 200: COMMON

// =============================================================
// swagger:route PUT /api/v1/admin/jwt/change-password v1-admin-jwt admin_change_password
// 管理员修改密码
// 在非法请求发出时（管理员不处于登入状态）会返回错误信息。
// PS: 通过postform形式传递密码, 不要使用url传参。
// responses:
// 200: COMMON
