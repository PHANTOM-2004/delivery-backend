package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/middleware/jwt"
	"delivery-backend/models"
	"delivery-backend/service/cache"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// merchant获得商家对应商店的信息, 返回所有的商铺
func GetRestaurants(c *gin.Context) {
	merchant_id := jwt.NewJwtInfo(c).GetID()
	res, err := models.GetRestaurantByMerchant(merchant_id)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	data := map[string]any{"restaurants": res}
	app.Response(c, http.StatusOK, ecode.SUCCESS, data)
}

const (
	// 店铺的状态放在redis中，而不是Table中
	RestaurantOpenStatus  = "1"
	RestaurantCloseStatus = "0"
)

// merchant设置商店的开闭状态,通过redis进行开闭设置
// restaurant的id在url参数之中，status在url参数之中
func SetRestaurantStatus(c *gin.Context) {
	// 如果出现参数不合法
	status := c.Param("status")
	if status != RestaurantOpenStatus &&
		status != RestaurantCloseStatus {
		app.ResponseInvalidParams(c)
		return
	}
	restaurant_id := c.GetUint("restaurant_id")

	// 设置status
	r := cache.NewRedisRestaurantStatus(restaurant_id)
	err := r.Set(status)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	app.ResponseSuccess(c)
}

// merchant设置商店的开闭状态,通过redis进行开闭设置/获得
// restaurant的id在url参数之中，status在url参数之中
// 这个接口用于请求商家设置的开闭状态；
func GetRestaurantStatus(c *gin.Context) {
	restaurant_id := c.GetUint("restaurant_id")
	r := cache.NewRedisRestaurantStatus(restaurant_id)
	status, err := r.Get()
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	data := map[string]any{"status": status}
	app.Response(c, http.StatusOK, ecode.SUCCESS, data)
}

// 应保证经过handler判定对应restaurant存在；
func UpdateRestaurant(c *gin.Context) {
	restaurant_id := c.GetUint("restaurant_id")
	var data models.Restaurant
	err := c.Bind(&data)
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return
	}

	err = app.ValidateStruct(data)
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return
	}

	err = models.UpdateRestaurant(restaurant_id, data)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	app.ResponseSuccess(c)
}

// 该函数判定店铺名是否重复，如果重复则不允许创建
// 创建时后端只校验最大长度；前端负责校验其他
func CreateRestaurant(c *gin.Context) {
	merchant_id := jwt.NewJwtInfo(c).GetID()

	var data models.Restaurant
	err := c.Bind(&data)
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return
	}
  //注意添加外键
	data.MerchantID = merchant_id
	err = app.ValidateStruct(data)
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return
	}

	// 如果店铺存在那么就不应该再次创建
	exist, err := models.ExistRestaurant(data.RestaurantName)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	if exist {
		app.Response(c, http.StatusOK, ecode.ERROR_RESTAURANT_EXIST, nil)
		return
	}

	err = models.CreateRestaurant(&data)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	app.ResponseSuccess(c)
}