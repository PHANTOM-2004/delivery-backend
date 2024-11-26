package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/models"

	"github.com/gin-gonic/gin"
)

func GetDishes(c *gin.Context) {
	var err error

	restaurant_id := c.GetUint("restaurant_id")
	dishes, err := models.GetDishes(restaurant_id)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	data := map[string]any{
		"dishes": dishes,
	}
	app.ResponseSuccessWithData(c, data)
}
