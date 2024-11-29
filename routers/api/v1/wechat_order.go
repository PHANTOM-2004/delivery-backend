package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/middleware/wechat"
	"delivery-backend/models"
	wechat_service "delivery-backend/service/wechat"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type CreateOrderRequest struct {
	Address      string `json:"address" validate:"max=100"`
	CustomerName string `json:"customer" validate:"max=20"`
	PhoneNumber  string `json:"phone_number" validate:"e164"`
}

func CreateOrder(c *gin.Context) {
	restaurant_id, err := strconv.Atoi(c.Param("restaurant_id"))
	if err != nil {
		app.ResponseInvalidParams(c)
		return
	}
	/////////////////////////////////////
	// parse请求
	var req CreateOrderRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		app.ResponseInvalidParams(c)
		log.Debug(err)
		return
	}
	err = app.ValidateStruct(&req)
	if err != nil {
		app.ResponseInvalidParams(c)
		log.Debug(err)
		return
	}

	/////////////////////////////////////
	//从session中验证购物车
	session := wechat.DefaultSession(c)
	cart, err := session.GetCart(uint(restaurant_id))
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	if len(cart) == 0 {
		// 购物车为空，不可能创建
		app.Response(c, http.StatusOK, ecode.ERROR_WX_CART_EMPTY, nil)
		return
	}
	info, err := session.GetInfo()
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	/////////////////////////////////////
	// NOTE:
	// 0.生成PickupNo
	// 1.生成OrderNo
	// 2.填入Address (from json)
	// 3.填入CustomerName (from json)
	// 4.绑定WechatUserID (from session)
	// 5.填入PhoneNumber(from json)
	order_gen := wechat_service.NewOrderGen(uint(restaurant_id))
	var order_no, pickup_number string
	pickup_number, err = order_gen.GetPickupNo()
	if err == nil {
		order_no, err = order_gen.GetOrderNo()
	}
	if err != nil {
		// 下单失败
		app.Response(c,
			http.StatusOK,
			ecode.ERROR_WX_ORDER_CREATE, nil)
		return
	}

	/////////////////////////////////////
	// 成功生成pickup no, order no
	order := models.Order{
		PickupNo:     pickup_number,
		OrderNo:      order_no,
		Address:      req.Address,
		PhoneNumber:  req.PhoneNumber,
		CustomerName: req.CustomerName,
		WechatUserID: info.ID,
	}
	err = models.CreateOrder(uint(restaurant_id), &order, cart)
	if err != nil {
		// 下单失败
		app.Response(c,
			http.StatusOK,
			ecode.ERROR_WX_ORDER_CREATE, nil)
		return
	}
	// 首先create order, 其次create order details
	app.ResponseSuccess(c)
}
