package docs

import (
	"bytes"
	"delivery-backend/models"
)

// swagger:parameters get_merchant_application
type GetMerchantApplication struct {
	// 页号
	// in: path
	Page int `json:"page"`
}

// https://github.com/go-swagger/go-swagger/blob/master/fixtures/goparsing/classification/operations/noparams.go#L28-L33
// swagger:parameters send_merchant_application
type SendMerchantApplication struct {
	// multipart/form data, 上传一个文件; 接受格式.png,.jpeg,.jpg
	// in: formData
	// required: true
	// swagger:file
	License *bytes.Buffer `json:"license"`
	// 最大长度: 50 bytes
	// required: true
	// in: formData
	Email string `json:"email"`
	// 手机格式为E.164规范
	// required: true
	// in: formData
	PhoneNumber string `json:"phone_number"`
	// 最大长度: 300 bytes
	// required: true
	// in: formData
	Description string `json:"description"`
	// 最大长度: 20 bytes
	// required: true
	// in: formData
	Name string `json:"name"`
}

// swagger:response merchant_application
type MerchantApplicationResponse struct {
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
			Applications []models.MerchantApplication `json:"applications"`
		} `json:"data"`
	}
}

//=============================================================
// swagger:route GET /api/v1/admin/jwt/merchant-application/{page} v1-admin get_merchant_application
// 管理员请求获得申请表, 每一页返回10个条目
// responses:
// 200: COMMON
// 200: merchant_application

//=============================================================
// swagger:route POST /api/v1/customer/merchant-application v1-customer send_merchant_application
// 顾客发起商务合作申请，上传对应信息以及执照
// responses:
// 200: COMMON
