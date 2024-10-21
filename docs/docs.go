// Package classification 父亲模拟器
//
// Documentation of our awesome API.
//
//	 Schemes: https
//	 BasePath: /
//	 Version: 0.1.0
//	 Host: some-url.com
//
//	 Consumes:
//	 - https
//
//	 Produces:
//	 - application/json
//
//	 Security:
//	 - basic
//
//	SecurityDefinitions:
//	basic:
//	  type: basic
//
// swagger:meta
package docs

// 注意：返回一个json数据
// 通常返回code=200, 请求被正确处理，单数需要校验msg字段是否为"ok";
// 如果"msg"不是"ok"则说明失败，msg含有对应报错信息.
// 具体返回的数据存在data中，详细见各个接口对于返回值的说明
// PS:如果出现服务器内部错误的code(500)，及时告诉开发者
// swagger:response COMMON
type CommonResponse struct {
	// in:body
	Body struct {
		// Required:true
		// Example: 200
		Code int `json:"code"`
		// Example: ok
		// error message
		Msg string `json:"msg"`
		// Required:true
		// data to get
		Data map[string]any `json:"data"`
	}
}
