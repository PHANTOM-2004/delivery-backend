package docs

// =============================================================
// 用于登入与认证的参数，注意后端处理的最大范围, 否则会返回错误
// swagger:parameters admin_auth admin_login admin_create admin_delete
type AdminAccountParam struct {
	// 账户名：最小长度10，最大长度30
	// required: true
	Account string `json:"account"`
}

// swagger:parameters admin_auth admin_login admin_create admin_change_password
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

// swagger:parameters  admin_change_password
type AccessToken struct {
	// 部分api调用的token
	// required: true
	AccessToken string `json:"access_token"`
}

//=============================================================
// swagger:route GET /admin/auth admin admin_auth
// JWT, 请求获得access_token
// 通过url参数发送请求; 获得返回返回字段在data中，key="access_token"
// responses:
// 200: COMMON

//=============================================================
// swagger:route POST /admin/login admin admin_login
// 登入的身份认证
// PS: 通过postform发送参数, 否则会认证错误
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
// swagger:route POST /api/v1/admin/change-password v1 admin_change_password
// 管理员修改密码
// 在非法请求发出时（管理员不处于登入状态）会返回错误信息。
// PS: 调用api时，access_token + session 双重认证
// responses:
// 200: COMMON
