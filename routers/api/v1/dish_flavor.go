package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/setting"
	"delivery-backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// 返回上传的名字
func uploadDishImage(c *gin.Context) (string, bool) {
	//////////////////// 校验文件
	file, err := c.FormFile("image")
	if err != nil {
		app.ResponseInvalidParams(c)
		return "", false
	}
	ext, v := setting.AppSetting.CheckDishImageExt(file.Filename)
	if !v {
		log.Debugf("wrong ext[%s]", ext)
		app.ResponseInvalidParams(c)
		return "", false
	}
	log.Debug("uploaded: ", file.Filename)
	// 重命名文件
	name := setting.AppSetting.GenDishImageName() + ext
	dst := setting.AppSetting.GetDishImageStorePath(name)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		app.ResponseInternalError(c, err)
		return "", false
	}
	log.Trace("saved to:", dst)
	return name, true
}

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
	name, v := uploadDishImage(c)
	if !v {
		return
	}
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
