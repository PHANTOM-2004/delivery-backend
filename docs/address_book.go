package docs

import (
	"delivery-backend/models"
	v1 "delivery-backend/routers/api/v1"
)

// swagger:response  customer_get_address_book
type GetAddressBookResponse struct {
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
			Cart []models.AddressBook `json:"address_books"`
		} `json:"data"`
	}
}

// =============================================================
// swagger:route GET /api/v1/wx/customer/addressbook v1-wechat customer_get_address_book
// 获得一个顾客的地簿
// 注意直接传递一个json即可
// 前端注意校验数据格式
// responses:
// 200: customer_get_address_book

//swagger:parameters customer_create_address_book
type CreateAddressBookRequest struct {
	//in:body
	Object struct {
		v1.CreateAddressBookRequest
	} `json:"object"`
}

// =============================================================
// swagger:route POST /api/v1/wx/customer/addressbook v1-wechat customer_create_address_book
// 顾客创建地址簿
// 前端注意校验数据格式
// responses:
// 200: COMMON

//swagger:parameters customer_update_address_book
type UpdateAddressBookRequest struct {
	//in:path
	ID uint `json:"address_book_id"`
	//in:body
	Object struct {
		v1.UpdateAddressBookRequest
	} `json:"object"`
}

// =============================================================
// swagger:route PUT /api/v1/wx/customer/addressbook/{address_book_id} v1-wechat customer_update_address_book
// 顾客更新地址簿，注意字段不必是必填
// 前端注意校验数据格式
// responses:
// 200: COMMON

//swagger:parameters customer_default_address_book
type SetDefaultAddressBookRequest struct {
	//in:path
	ID uint `json:"address_book_id"`
}

// =============================================================
// swagger:route PUT /api/v1/wx/customer/addressbook/{address_book_id}/default v1-wechat customer_default_address_book
// 顾客更新地址簿，注意字段不必是必填
// 前端注意校验数据格式
// responses:
// 200: COMMON

//swagger:parameters customer_delete_address_book
type DeleteAddressBookRequest struct {
	//in:path
	ID uint `json:"address_book_id"`
}

// =============================================================
// swagger:route DELETE /api/v1/wx/customer/addressbook/{address_book_id} v1-wechat customer_delete_address_book
// 顾客更新地址簿，注意字段不必是必填
// 前端注意校验数据格式
// responses:
// 200: COMMON
