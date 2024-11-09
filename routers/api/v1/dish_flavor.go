package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/internal/setting"
	"delivery-backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// 商家create一个dish, dish属于某个商铺
// 传递参数
// 1. restaurant
// 2. category
// 3. form
func CreateDish(c *gin.Context) {
	category_id, err := strconv.Atoi(c.Param("category_id"))
	if err != nil {
		app.ResponseInvalidParams(c)
		return
	}

	//////////////////// 校验是否是自己的dish
	restaurant_id := c.GetUint("restaurant_id")
	category, err := models.GetCategory(uint(category_id))
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	} else if category.ID == 0 {
		// 找不到category
		app.Response(c, http.StatusOK, ecode.ERROR_CATEGORY_NOT_FOUND, nil)
		return
	} else if category.RestaurantID != restaurant_id {
		// 不是自己的category
		app.Response(c, http.StatusOK, ecode.ERROR_MERCHANT_UNAUTH, nil)
		return
	}

	//////////////////// 校验form
	var dish models.Dish
	err = c.Bind(&dish)
	if err != nil {
		app.ResponseInvalidParams(c)
		return
	}
	err = app.ValidateStruct(dish)
	if err != nil {
		app.ResponseInvalidParams(c)
		return
	}
	//////////////////// 校验文件
	file, err := c.FormFile("image")
	if err != nil {
		app.ResponseInvalidParams(c)
		return
	}
	ext, v := setting.AppSetting.CheckDishImageExt(file.Filename)
	if !v {
		log.Debugf("wrong ext[%s]", ext)
		app.ResponseInvalidParams(c)
		return
	}
	log.Debug("uploaded: ", file.Filename)
	// 重命名文件
	name := setting.AppSetting.GenDishImageName() + ext
	dst := setting.AppSetting.GetDishImageStorePath(name)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	log.Trace("saved to:", dst)

	// 更新待上传的文件
	dish.Image = name
	// 加入category_id
	dish.CategoryID = uint(category_id)

	err = models.CreateDish(&dish)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccess(c)
}

// 注意校验dish是否属于对应的category
// 在中间件中已经校验过restaurant的所属
// 传递参数
// 1. restaurant_id
// 2. category_id
// 3. dish_id
func UpdateDish(c *gin.Context) {
	// TODO:更改图片
}

func GetDishes(c *gin.Context) {
}
