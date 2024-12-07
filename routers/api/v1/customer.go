package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type MerchantApplyReq struct {
	Description string `json:"description" validate:"min=1,max=300"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required,e164"`
	Name        string `json:"name" validate:"min=2,max=20"`
	Path        string `json:"path" validate:"max=80"`
}

// 由顾客/骑手发起申请，提出商务合作请求
// https://gin-gonic.com/docs/examples/upload-file/single-file/
func MerchantApply(c *gin.Context) {
	req := MerchantApplyReq{}
	err := c.ShouldBindJSON(&req)
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
	data := models.MerchantApplication{
		Description: req.Description,
		License:     req.Path,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
	}

	err = models.CreateMerchantApplication(&data)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccess(c)
}

func GetRestaurantCategoryDish(c *gin.Context) {
	var err error
	restaurant_id, err := strconv.Atoi(c.Param("restaurant_id"))
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return
	}

	res, err := models.GetCategoryDishFlavor(uint(restaurant_id))
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccessWithData(c, map[string]any{"categories": res})
}
