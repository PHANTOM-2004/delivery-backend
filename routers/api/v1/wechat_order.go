package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/middleware/wechat"
	"delivery-backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	restaurant_id, err := strconv.Atoi(c.Param("restaurant_id"))
	if err != nil {
		app.ResponseInvalidParams(c)
		return
	}

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

	order := models.Order{}
	// TODO:
	// 0.生成PickupNo
	// 1.生成OrderNo
	// 2.填入Address (from json)
	// 3.填入CustomerName (from json)
	// 4.绑定WechatUserID (from session)
	// 5.填入PhoneNumber(from json)

	err = models.CreateOrder(uint(restaurant_id), &order, cart)
	if err != nil {
		// 下单失败
		app.Response(c, http.StatusOK, ecode.ERROR_WX_ORDER_CREATE, nil)
		return
	}
	// 首先create order, 其次create order details
	app.ResponseSuccess(c)
}
