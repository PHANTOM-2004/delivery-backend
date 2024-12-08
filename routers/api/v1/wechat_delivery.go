package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/middleware/wechat"
	"delivery-backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type RiderInfoReq struct {
	StuName string `json:"student_name" validate:"max=20"`
	StuNo   string `json:"student_no" validate:"max=7"`
	StuCard string `json:"student_card"`
}

func UploadRiderApplication(c *gin.Context) {
	req := RiderInfoReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.ResponseInvalidParams(c)
		return
	}
	err = app.ValidateStruct(&req)
	if err != nil {
		app.ResponseInvalidParams(c)
		return
	}

	session := wechat.DefaultSession(c)
	info, err := session.GetInfo()
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	err = models.CreateRiderApplication(&models.RiderApplication{
		WechatUserID: info.ID,
		StudentName:  req.StuName,
		StudentNo:    req.StuNo,
		StudentCard:  req.StuCard,
	})
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccess(c)
}

func AuthRider(c *gin.Context) bool {
	session := wechat.DefaultSession(c)
	info, err := session.GetInfo()
	if err != nil {
		app.ResponseInternalError(c, err)
		return false
	}
	if info.Role != models.RoleRider {
		log.Warn("non rider request rider api")
		app.Response(c, http.StatusUnauthorized, ecode.ERROR, nil)
		return false
	}
	return true
}

func GetDeliveryOrder(c *gin.Context) {
	status, err := strconv.Atoi(c.Param("status"))
	if err != nil {
		app.ResponseInvalidParams(c)
		return
	}
	if status < models.OrderNotPayed || status > models.OrderCanceled {
		app.ResponseInvalidParams(c)
		return
	}
	if !AuthRider(c) {
		return
	}

	// 获取对应status状态的
	orders, err := models.GetOrderByStatus(uint8(status))
	for i := range orders {
		// 填充骑手的response
		restaurant := orders[i].Restaurant
		orders[i].RestaurantInfoEx = &models.RestaurantInfoEx{
			Name:    restaurant.RestaurantName,
			Address: restaurant.Address,
		}
	}

	app.ResponseSuccessWithData(c, map[string]any{
		"orders": orders,
	})
}

func SetDeliveryOrderStatus(c *gin.Context) {
	status, err := strconv.Atoi(c.Param("status"))
	if err != nil {
		app.ResponseInvalidParams(c)
		return
	}
	order_id, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		app.ResponseInvalidParams(c)
		return
	}
	if status < models.OrderToDeliver || status > models.OrderFinished {
		app.ResponseInvalidParams(c)
		return
	}
	if !AuthRider(c) {
		return
	}
	success, err := models.SetOrderStatus(uint(order_id), uint8(status))
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccessWithData(c, map[string]any{
		"success": success,
	})
}
