package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/middleware/wechat"
	"delivery-backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// 获得一个商店的comment
func WXGetRestaurantComments(c *gin.Context) {
	restaurant_id, err := strconv.Atoi(c.Param("restaurant_id"))
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return
	}

	res, err := models.GetCommentsByRestaurnat(uint(restaurant_id))
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccessWithData(c, map[string]any{
		"comments": res,
	})
}

type CreateCommentRequest struct {
	Images  []string `json:"images" validate:"dive,max=100"`
	Content string   `json:"content" validate:"max=300"`
	Rating  uint8    `json:"rating" validate:"gte=0,lte=10"`
	OrderID uint     `json:"order_id" validate:"gte=1"`
}

func WXCreateComment(c *gin.Context) {
	restaurant_id, err := strconv.Atoi(c.Param("restaurant_id"))
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return
	}
	req := CreateCommentRequest{}
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

	comment := models.Comment{
		Rating:       req.Rating,
		Content:      req.Content,
		OrderID:      req.OrderID,
		WechatUserID: info.ID,
		RestaurantID: uint(restaurant_id),
	}

	err = models.CreateComment(&comment, req.Images)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccess(c)
}
