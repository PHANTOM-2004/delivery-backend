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

func CancelOrder(c *gin.Context) {
	order_id, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return
	}

	success, err := models.CancelOrder(uint(order_id))
	if err != nil || !success {
		log.Error("error cancel order")
		app.Response(c, http.StatusOK, ecode.ERROR_WX_ORDER_CANCEL, nil)
		return
	}
	app.ResponseSuccess(c)
}

func GetCustomerOrders(c *gin.Context) {
	session := wechat.DefaultSession(c)
	info, err := session.GetInfo()
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	res, err := models.GetOrderByUserID(info.ID)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccessWithData(c,
		map[string]any{
			"orders": res,
		})
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
	log.Tracef("user:[%v] is creating order", info.ID)

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
	log.Tracef("pickup_no:[%s], order_no:[%s]", pickup_number, order_no)

	/////////////////////////////////////
	// 成功生成pickup no, order no
	order := models.Order{
		PickupNo:     pickup_number,
		OrderNo:      order_no,
		Address:      req.Address,
		PhoneNumber:  req.PhoneNumber,
		CustomerName: req.CustomerName,
		WechatUserID: info.ID,
		RestaurantID: uint(restaurant_id),
	}
	err = models.CreateOrder(&order, cart)
	if err != nil {
		// 下单失败
		log.Debugf("user:[%v] fail to create order", info.ID)
		app.Response(c,
			http.StatusOK,
			ecode.ERROR_WX_ORDER_CREATE, nil)
		return
	}
	log.Trace("new order created with id", order.ID)
	// 仍需要清空购物车
	err = session.UpdateCart(uint(restaurant_id), []wechat.WXSessionCartStore{})
	if err != nil {
		app.Response(c,
			http.StatusOK,
			ecode.ERROR_WX_ORDER_CREATE, nil)
		log.Error("error cleaning carts")
		return
	}
	log.Trace("clear restaurant carts")

	// 首先create order, 其次create order details
	// 返回order_id
	app.ResponseSuccessWithData(c, map[string]any{
		"order_id": order.ID,
	})
}
