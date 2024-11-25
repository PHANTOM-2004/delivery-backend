package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func AddCategoryDish(c *gin.Context) {
	category_id, err := strconv.ParseUint(c.Param("category_id"), 10, 0)
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return
	}

	dishes_id := app.NewIDArrayParser("dishes", c).Parse()
	if len(dishes_id) == 0 {
		app.ResponseInvalidParams(c)
		return
	}

	log.Debugf("add dishes to category[%v]", category_id)
	log.Debug("add dishes:", dishes_id)
	err = models.AddCategoryDish(uint(category_id), dishes_id)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccess(c)
}

func DeleteCategoryDish(c *gin.Context) {
	category_id, err := strconv.ParseUint(c.Param("category_id"), 10, 0)
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return
	}

	dishes_id := app.NewIDArrayParser("dishes", c).Parse()
	if len(dishes_id) == 0 {
		app.ResponseInvalidParams(c)
		return
	}

	log.Debugf("delete dishes from category[%v]", category_id)
	log.Debug("delete dishes:", dishes_id)
	err = models.DeleteCategoryDish(uint(category_id), dishes_id)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccess(c)
}
