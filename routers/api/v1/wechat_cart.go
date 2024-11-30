package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/middleware/wechat"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type WXUpdateCartRequest struct {
	// in: body
	Cart []wechat.WXSessionCartStore `json:"cart"`
}

func WXUpdateCart(c *gin.Context) {
	restaurant_id, err := strconv.Atoi(c.Param("restaurant_id"))
	if err != nil {
		app.ResponseInvalidParams(c)
		log.Debug(err)
		return
	}
	cart := WXUpdateCartRequest{}
	err = c.ShouldBindJSON(&cart)
	if err != nil {
		app.ResponseInvalidParams(c)
		log.Debug(err)
		return
	}
	session := wechat.DefaultSession(c)
	err = session.UpdateCart(uint(restaurant_id), cart.Cart)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccess(c)
}

func WXGetCart(c *gin.Context) {
	restaurant_id, err := strconv.Atoi(c.Param("restaurant_id"))
	if err != nil {
		app.ResponseInvalidParams(c)
		log.Debug(err)
		return
	}
	session := wechat.DefaultSession(c)
	cart, err := session.GetCart(uint(restaurant_id))
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	res := map[string]any{
		"cart": cart,
	}
	app.ResponseSuccessWithData(c, res)
}
