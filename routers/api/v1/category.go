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

type categoryRequest struct {
	Name string `form:"name" validate:"max=30"`
	// 分类类型，1代表菜品，2代表套餐;默认是菜品
	Type uint8 `form:"type" validate:"gte=1,lte=2"`
	// 排序值，用于决定顺序；
	Sort uint16 `form:"sort" validate:"gte=0"`
	// 1代表禁用，2代表启用，默认禁用
	Status uint8 `form:"status" validate:"gte=1,lte=2"`
}

func (r *categoryRequest) GetCategoryModel() *models.Category {
	category := models.Category{
		Name:   r.Name,
		Type:   r.Type,
		Sort:   r.Sort,
		Status: r.Status,
	}
	return &category
}

// 暂不去处理重复的插入类别
func CreateCategory(c *gin.Context) {
	restaurant_id := c.GetUint("restaurant_id")
	var data categoryRequest
	// 使用bind的时候会返回http.StatusBadrequest
	err := c.Bind(&data)
	// 姓名不能空
	if err != nil || data.Name == "" {
		log.Warn(err)
		app.ResponseInvalidParams(c)
		return
	}

	err = app.ValidateStruct(&data)
	if err != nil {
		log.Warn(err)
		app.ResponseInvalidParams(c)
		return
	}

	category := data.GetCategoryModel()
	// 设置外键
	category.RestaurantID = restaurant_id

	err = models.CreateCategory(category)
	if err != nil {
		log.Warn(err)
		app.ResponseInternalError(c, err)
		return
	}

	app.ResponseSuccess(c)
}

// 需要鉴权, 修改的是不是商家自己的
// category_id作为url参数
// 更新某一个Category; 现在
func UpdateCategory(c *gin.Context) {
	category_id, _ := strconv.Atoi(c.Param("category_id"))
	// 使用bind的时候会返回http.StatusBadrequest
	var data categoryRequest
	err := c.Bind(&data)
	if err != nil {
		log.Warn(err)
		app.ResponseInvalidParams(c)
		return
	}
	err = app.ValidateStruct(&data)
	if err != nil {
		log.Warn(err)
		app.ResponseInvalidParams(c)
		return
	}

	category := data.GetCategoryModel()
	err = models.UpdateCategory(uint(category_id), category)
	// 记录必然找到
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccess(c)
}

// category_id作为url参数
func DeleteCategory(c *gin.Context) {
	category_id, _ := strconv.Atoi(c.Param("category_id"))

	err := models.DeleteCategory(uint(category_id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		app.Response(c, http.StatusOK, ecode.ERROR_CATEGORY_NOT_FOUND, nil)
		return
	}
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccess(c)
}
