package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func WXGetRestaurant(c *gin.Context) {
	restaurant_id, err := strconv.Atoi(c.Param("restaurant_id"))
	if err != nil {
		app.ResponseInvalidParams(c)
		return
	}
	restaurants, err := models.GetRestaurantByID(uint(restaurant_id))
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccessWithData(c, map[string]any{
		"address": restaurants.Address,
		"name":    restaurants.RestaurantName,
	})
}
