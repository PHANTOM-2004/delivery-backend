package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// TODO:设置一个恰当的脚本进行整个流程的参数验证
// 可以给数据库中插入几个已知的商家数据,用于测试的时候使用

// 获得某一家店铺的所有的Category
// restaurant_id作为url参数
func GetCategories(c *gin.Context) {
	restaurant_id := c.GetUint("restaurant_id")
	categories, err := models.GetCategoryByRestaurant(uint(restaurant_id))
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	data := map[string]any{"categories": categories}
	app.Response(c, http.StatusOK, ecode.SUCCESS, data)
}

// 暂不去处理重复的插入类别
func CreateCategory(c *gin.Context) {
	restaurant_id := c.GetUint("restaurant_id")
	var data models.Category
	// 使用bind的时候会返回http.StatusBadrequest
	err := c.Bind(&data)
	// 姓名不能空
	if err != nil || data.Name == "" {
		log.Warn(err)
		app.ResponseInvalidParams(c)
		return
	}

	err = app.ValidateStruct(data)
	if err != nil {
		log.Warn(err)
		app.ResponseInvalidParams(c)
		return
	}

	// 设置外键
	data.RestaurantID = restaurant_id

	err = models.CreateCategory(&data)
	if err != nil {
		log.Warn(err)
		app.ResponseInternalError(c, err)
		return
	}

	app.ResponseSuccess(c)
}

// restaurant_id作为url参数
// category_id作为url参数
// 更新某一个Category; 现在
func UpdateCategory(c *gin.Context) {
	restaurant_id := c.GetUint("restaurant_id")
	category_id, err := strconv.Atoi(c.Param("category_id"))
	if err != nil {
		app.ResponseInvalidParams(c)
		return
	}

	var data models.Category

	// 使用bind的时候会返回http.StatusBadrequest
	err = c.Bind(&data)
	if err != nil {
		log.Warn(err)
		app.ResponseInvalidParams(c)
		return
	}

	err = app.ValidateStruct(data)
	if err != nil {
		log.Warn(err)
		app.ResponseInvalidParams(c)
		return
	}

	// NOTE:此处由于验证了商家必然更改自己的商铺，
	// 但是也需要验证category对应的外键是否是restaurant_id
	category, err := models.GetCategory(uint(category_id))
	if category.RestaurantID != restaurant_id {
		// 商家修改的不是自己店铺的category
		app.Response(c, http.StatusUnauthorized, ecode.ERROR_MERCHANT_UNAUTH, nil)
		return
	}

	err = models.UpdateCategory(uint(category_id), data)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		app.Response(c, http.StatusOK, ecode.ERROR_CATEGORY_NOT_FOUND, nil)
		return
	} else if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	app.ResponseSuccess(c)
}
