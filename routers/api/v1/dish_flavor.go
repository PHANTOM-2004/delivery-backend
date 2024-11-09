package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/internal/setting"
	"delivery-backend/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// 如果中间有校验失败，会返回false
// 校验成功时，返回上传的image名字
func uploadDishImage(c *gin.Context) (string, bool) {
	//////////////////// 校验文件
	file, err := c.FormFile("image")
	if err != nil {
		log.Debug(err)
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

// 如果中间有校验失败，会返回false
// 校验成功时，返回上传的dish
func validateDish(c *gin.Context) (*models.Dish, bool) {
	//////////////////// 校验form
	var dish models.Dish
	err := c.Bind(&dish)
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return nil, false
	}
	err = app.ValidateStruct(dish)
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return nil, false
	}

	return &dish, true
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
	dish, v := validateDish(c)
	if !v {
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

	err = models.CreateDish(dish)
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
// TODO: 后续可以做一个update hook用于清理之前的图片
func UpdateDish(c *gin.Context) {
	category_id, err := strconv.Atoi(c.Param("category_id"))
	if err != nil {
		app.ResponseInvalidParams(c)
		return
	}
	dish_id, err := strconv.Atoi(c.Param("dish_id"))
	if err != nil {
		app.ResponseInvalidParams(c)
		return
	}

	//////////////////// 校验form
	dish, v := validateDish(c)
	if !v {
		return
	}

	/////验证dish是否属于当前category
	d, err := models.GetDish(uint(dish_id))
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	if d.ID == 0 {
		// 没有找到dish
		app.Response(c, http.StatusOK, ecode.ERROR_DISH_NOT_FOUND, nil)
		return

	}
	if d.CategoryID != uint(category_id) {
		// 无权修改他人的dish
		app.Response(c, http.StatusUnauthorized, ecode.ERROR_MERCHANT_UNAUTH, nil)
		return
	}

	//////////////////// 校验文件
	_, err = c.FormFile("image")
	if !errors.Is(err, http.ErrMissingFile) {
		// 有文件上传的时候才考虑更新文件
		name, v := uploadDishImage(c)
		if !v {
			return
		}
		dish.Image = name
	}

	/////更新dish, 此时记录必然存在
	err = models.UpdateDish(uint(dish_id), *dish)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	app.ResponseSuccess(c)
}
