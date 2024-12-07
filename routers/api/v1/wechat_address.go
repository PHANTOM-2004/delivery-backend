package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/middleware/wechat"
	"delivery-backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type UpdateAddressBookRequest struct {
	// 最大长度80
	Address string `json:"address" validate:"max=80"`
	// 最大长度20
	Name string `json:"name" validate:"max=20"`
	// 1或者2,1表示男性，2表示女性
	Gender uint8 `json:"gender" validate:"gte=1,lte=2"`
	// e164格式
	PhoneNumber string `json:"phone_number" validate:"e164"`
}

func (r *UpdateAddressBookRequest) GetModel() *models.AddressBook {
	return &models.AddressBook{
		Address:     r.Address,
		Name:        r.Name,
		Gender:      r.Gender,
		PhoneNumber: r.PhoneNumber,
	}
}

func UpdateAddressBook(c *gin.Context) {
	address_book_id, err := strconv.Atoi(c.Param("address_book_id"))
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return
	}
	req := UpdateAddressBookRequest{}
	err = c.ShouldBindJSON(&req)
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return
	}
	err = app.ValidateStruct(&req)
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return
	}
	err = models.UpdateAddressBook(uint(address_book_id), req.GetModel())
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccess(c)
}

func SetDefaultAddressBook(c *gin.Context) {
	// TODO:
	address_book_id, err := strconv.Atoi(c.Param("address_book_id"))
	if err != nil {
		app.ResponseInvalidParams(c)
		return
	}

	info, err := wechat.DefaultSession(c).GetInfo()
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	err = models.SetDefaultAddressBook(info.ID, uint(address_book_id))
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccess(c)
}

type CreateAddressBookRequest struct {
	// 最大长度80
	// required: true
	Address string `json:"address" validate:"max=80,required"`
	// 最大长度20
	// required: true
	Name string `json:"name" validate:"max=20,required"`
	// 1或者2,前者是男性
	// required: true
	Gender uint8 `json:"gender" validate:"gte=1,lte=2,required"`
	// e164格式
	// required: true
	PhoneNumber string `json:"phone_number" validate:"e164,required"`
}

func (r *CreateAddressBookRequest) GetModel() *models.AddressBook {
	return &models.AddressBook{
		Address:     r.Address,
		Name:        r.Name,
		Gender:      r.Gender,
		PhoneNumber: r.PhoneNumber,
	}
}

func CreateAddressBook(c *gin.Context) {
	var err error
	req := CreateAddressBookRequest{}
	err = c.ShouldBindJSON(&req)
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return
	}
	err = app.ValidateStruct(&req)
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return
	}
	info, err := wechat.DefaultSession(c).GetInfo()
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	m := req.GetModel()
	m.WechatUserID = info.ID
	err = models.CreateAddressBook(m)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccess(c)
}

func DeleteAddressBook(c *gin.Context) {
	address_book_id, err := strconv.Atoi(c.Param("address_book_id"))
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return
	}
	err = models.DeleteAddressBook(uint(address_book_id))
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccess(c)
}

func GetAddressBook(c *gin.Context) {
	session := wechat.DefaultSession(c)
	w, err := session.GetInfo()
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	res, err := models.GetAddressBooks(w.ID)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccessWithData(
		c,
		map[string]any{
			"address_books": res,
		},
	)
}
