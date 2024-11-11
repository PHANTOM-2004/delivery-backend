package merchant_service

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/models"
	handler "delivery-backend/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func DishAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		dish_id, err := strconv.Atoi(c.Param("dish_id"))
		if err != nil {
			app.ResponseInvalidParams(c)
			return
		}
		merchant_id := handler.NewMerchInfoHanlder(c).GetID()
		d_merchant, err := models.GetMerchantByDish(uint(dish_id))
		if err != nil {
			app.ResponseInternalError(c, err)
			return
		}
		if d_merchant.ID != merchant_id {
			// 无权修改不是自己的
			log.Debugf("current merchant_id[%v]", merchant_id)
			log.Debugf("dish[%v] belongs to[%v]", dish_id, d_merchant.ID)
			app.Response(c, http.StatusUnauthorized, ecode.ERROR_MERCHANT_UNAUTH, nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

// 根据参数category_id, restaurant_id, 校验修改的category是否是商家自己的
// 避免api被滥用于越权修改
func CategoryAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		category_id, err := strconv.Atoi(c.Param("category_id"))
		if err != nil {
			app.ResponseInvalidParams(c)
			return
		}
		merchant_id := handler.NewMerchInfoHanlder(c).GetID()
		c_merchant, err := models.GetMerchantByCategory(uint(category_id))
		if err != nil {
			app.ResponseInternalError(c, err)
			return
		}
		if merchant_id != c_merchant.ID {
			// 无权修改不是自己的
			log.Debugf("current merchant_id[%v]", merchant_id)
			log.Debugf("category[%v] belongs to[%v]", category_id, c_merchant.ID)
			app.Response(c, http.StatusUnauthorized, ecode.ERROR_MERCHANT_UNAUTH, nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

// 对入参restaurant_id进行验证，验证是否存在这个restaurant
// 对参数进行验证，对商家身份进行验证，商家必须访问的是自己的店铺
// 如果通过验证，会把restaurant_id设置在gin.Context
func RestaurantAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 首先商家鉴权，
		merchant_id := handler.NewMerchInfoHanlder(c).GetID()
		restaurant_id, err := strconv.Atoi(c.Param("restaurant_id"))
		if err != nil {
			app.ResponseInvalidParams(c)
			return
		}

		r, err := models.GetRestaurantByID(uint(restaurant_id))
		if r.ID == 0 {
			// restaurant不存在
			app.Response(c, http.StatusOK, ecode.ERROR_RESTAURANT_NOT_FOUND, nil)
			c.Abort()
			return
		}

		// NOTE:这里做一次鉴权，考虑到商家修改店铺状态实际上是低频事件；
		// 我们并不希望商家修改的不是自己的店铺
		if r.MerchantID != merchant_id {
			app.Response(c, http.StatusUnauthorized, ecode.ERROR_MERCHANT_UNAUTH, nil)
			c.Abort()
			return
		}

		// 设置可靠的restaurant_id在上下文中
		c.Set("restaurant_id", r.ID)
		c.Next()
	}
}
