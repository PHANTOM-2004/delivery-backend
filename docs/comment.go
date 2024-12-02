package docs

import (
	"bytes"
	"delivery-backend/models"
	v1 "delivery-backend/routers/api/v1"

	"github.com/go-openapi/runtime"
)

// swagger:parameters customer_get_comment_image
type CustomerGetCommentImageRequest struct {
	// example: comment-02fd3ce1-fdcb-4a30-94a4-db3f9c241871.png
	// Required:true
	// in: path
	ID string `json:"*filepath"`
}

// swagger:response customer_get_comment_image
type CustomerGetCommentImageResponse struct {
	// Required:true
	Image runtime.File `json:"image"`
}

// =============================================================
// swagger:route GET /api/v1/wx/customer/comment/{*file_path} v1-wechat customer_get_comment_image
// 请求得到评论照片; 如果路径内容不存在则对应httpcode=404
// responses:
// 200: customer_get_comment_image

//swagger:parameters customer_upload_comment_image
type CustomerUploadImageRequest struct {
	// in: formData
	// required: true
	// swagger:file
	License *bytes.Buffer `json:"image"`
}

// =============================================================
// swagger:route POST /api/v1/wx/customer/comment/image v1-wechat customer_upload_comment_image
// 上传评论照片
// responses:
// 200: CustomerGetCommentImageResponse

// swagger:parameters customer_create_comment
type CustomerCreateCommentRequest struct {
	//in:path
	//required:true
	RestaurantID uint `json:"restaurant_id"`
	//in:body
	Body struct {
		v1.CreateCommentRequest
	}
}

// =============================================================
// swagger:route POST /api/v1/wx/customer/comment/restaurant/{restaurant_id} v1-wechat customer_create_comment
// 发布评论
// responses:
// 200: COMMON

// swagger:parameters customer_get_comments
type CustomerGetComments struct {
	//in:path
	//required:true
	RestaurantID uint `json:"restaurant_id"`
}

// swagger:response customer_get_comments
type CustomerGetCommentsResponse struct {
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
			Applications []models.Comment `json:"comments"`
		} `json:"data"`
	}
}

// =============================================================
// swagger:route GET /api/v1/wx/customer/comment/restaurant/{restaurant_id} v1-wechat customer_get_comments
// 获得店铺的评论
// responses:
// 200: customer_get_comments
