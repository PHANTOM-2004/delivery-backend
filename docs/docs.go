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
//	 - multipart/form-data
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

import "time"

// 注意：返回一个json数据,httpcode在请求头而不在返回的json数据中
// (1) 通常返回htppcode=200, 说明请求正确处理，不代表请求成功;
// 如果"msg"不是"ok"则说明失败(或者ecode!=200)，msg也含有对应报错信息.
// 具体返回的数据存在data中，详细见各个接口对于返回值的说明
// (2) 注意token鉴权失败的情况，需要刷新access_token
// (3) 需要判断业务逻辑码ecode, ecode=200是成功处理请求
//
// (*) 如果出现httpcode=500，及时告诉开发者
// swagger:response COMMON
type CommonResponse struct {
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
    // Example: {}
		Data map[string]any `json:"data"`
	}
}

// https://github.com/go-swagger/go-swagger/blob/master/fixtures/goparsing/classification/operations/noparams.go#L28-L33
//

type DefaultModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
