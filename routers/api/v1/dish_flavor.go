package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/internal/setting"
	"delivery-backend/models"
	handler "delivery-backend/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
func validateDish(c *gin.Context) (*dishRequest, bool) {
	//////////////////// 校验form
	var dish dishRequest
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

// 注意与models.Dish信息同步
type dishRequest struct {
	Name        string `form:"name" validate:"max=30"`
	Price       uint   `form:"price"`
	Description string `form:"description" validate:"max=50"`
	Sort        uint16 `form:"sort"`
}

func (r *dishRequest) GetDishModel() *models.Dish {
	dish := models.Dish{
		Price:       r.Price,
		Name:        r.Name,
		Sort:        r.Sort,
		Description: r.Description,
	}
	return &dish
}

// 商家create一个dish, dish属于某个商铺
// 传递参数
// 1. restaurant
// 2. category
// 3. form
func CreateDish(c *gin.Context) {
	category_id, _ := strconv.Atoi(c.Param("category_id"))
	var err error

	//////////////////// 校验form
	dish_r, v := validateDish(c)
	if !v {
		return
	}
	//////////////////// 校验文件
	image_name, v := uploadDishImage(c)
	if !v {
		return
	}

	dish := dish_r.GetDishModel()
	dish.Image = image_name
	dish.CategoryID = uint(category_id)
	err = models.CreateDish(dish)
	if err != nil {
		log.Debug(dish)
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccess(c)
}

func DeleteDish(c *gin.Context) {
	dish_id, _ := strconv.Atoi(c.Param("dish_id"))
	err := models.DeleteDish(uint(dish_id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		app.Response(c, http.StatusOK, ecode.ERROR_DISH_NOT_FOUND, nil)
		return
	}
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
	dish_id, _ := strconv.Atoi(c.Param("dish_id"))

	var err error
	//////////////////// 校验form
	dish_r, v := validateDish(c)
	if !v {
		return
	}

	dish := dish_r.GetDishModel()
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
	err = models.UpdateDish(uint(dish_id), dish)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	app.ResponseSuccess(c)
}

// 通过url path传参即可，传递name
func CreateFlavor(c *gin.Context) {
	restaurant_id := c.GetUint("restaurant_id")
	if restaurant_id == 0 {
		log.Warn("restaurant id could not be 0")
	}

	name := c.Param("name")
	if name == "" {
		log.Debug("name could not be empty")
		app.ResponseInvalidParams(c)
		return
	}

	f := models.Flavor{
		Name:         name,
		RestaurantID: restaurant_id,
	}
	err := models.CreateFlavor(&f)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccess(c)
}

func DeleteFlavor(c *gin.Context) {
	flavor_id, err := strconv.Atoi(c.Param("flavor_id"))
	if err != nil {
		app.ResponseInvalidParams(c)
		return
	}
	// 验证修改的是不是自己的
	merchant_id, err := models.GetMerchantIDByFlavor(uint(flavor_id))
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	if merchant_id != handler.NewMerchInfoHanlder(c).GetID() {
		app.Response(c, http.StatusUnauthorized, ecode.ERROR_MERCHANT_UNAUTH, nil)
		return
	}

	err = models.DeleteFlavor(uint(flavor_id))
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccess(c)
}

func UpdateFlavor(c *gin.Context) {
	flavor_id, err := strconv.Atoi(c.Param("flavor_id"))
	if err != nil {
		app.ResponseInvalidParams(c)
		return
	}
	// 验证修改的是不是自己的
	merchant_id, err := models.GetMerchantIDByFlavor(uint(flavor_id))
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	if merchant_id != handler.NewMerchInfoHanlder(c).GetID() {
		app.Response(c, http.StatusUnauthorized, ecode.ERROR_MERCHANT_UNAUTH, nil)
		return
	}

	name := c.Param("name")
	if name == "" {
		log.Debug("name could not be empty")
		app.ResponseInvalidParams(c)
		return
	}

	err = models.UpdateFlavor(uint(flavor_id), name)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccess(c)
}

func GetDishFlavor(c *gin.Context) {
	dish_id, err := strconv.Atoi(c.Param("dish_id"))
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return
	}
	d, err := models.GetDishFlavors(uint(dish_id))
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	res := map[string]any{
		"flavors": d.Flavors,
	}

	app.Response(c, http.StatusOK, ecode.SUCCESS, res)
}

func DeleteDishFlavor(c *gin.Context) {
	dish_id, err := strconv.ParseUint(c.Param("dish_id"), 10, 0)
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return
	}
	flavors_id := app.NewIDArrayParser("flavors", c).Parse()
	if len(flavors_id) == 0 {
		log.Debug(" flavors参数有误")
		app.ResponseInvalidParams(c)
		return
	}

	log.Debugf("delete flavors from dish[%v]", dish_id)
	log.Debug("delete flavors:", flavors_id)
	err = models.DeleteDishFlavor(uint(dish_id), flavors_id)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccess(c)
}

// 传入flavors作为flavor数组
func AddDishFlavor(c *gin.Context) {
	dish_id, err := strconv.ParseUint(c.Param("dish_id"), 10, 0)
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return
	}

	flavors_id := app.NewIDArrayParser("flavors", c).Parse()
	if len(flavors_id) == 0 {
		app.ResponseInvalidParams(c)
		return
	}

	log.Debugf("add flavors to dish[%v]", dish_id)
	log.Debug("add flavors:", flavors_id)
	err = models.AddDishFlavor(uint(dish_id), flavors_id)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccess(c)
}

func GetRestaurantFlavors(c *gin.Context) {
	restaurant_id := c.GetUint("restaurant_id")
	flavors, err := models.GetFlavors(restaurant_id)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	res := map[string]any{
		"flavors": flavors,
	}
	app.Response(c, http.StatusOK, ecode.SUCCESS, res)
}
