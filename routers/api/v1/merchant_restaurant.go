package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/models"
	handler "delivery-backend/service"
	"delivery-backend/service/cache"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// merchant获得商家对应商店的信息, 返回所有的商铺
func GetRestaurants(c *gin.Context) {
	merchant_id := handler.NewMerchInfoHanlder(c).GetID()
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

	log.Tracef("set restaurant[%v] status[%v]", restaurant_id, status)

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
	log.Tracef("get restaurant[%v] status[%v]", restaurant_id, status)
	data := map[string]any{"status": status}
	app.Response(c, http.StatusOK, ecode.SUCCESS, data)
}

type restaurantRequest struct {
	RestaurantName string `form:"restaurant_name" validate:"max=50"`
	// 店铺的地址
	Address string `validate:"max=50"`
	// 商铺简介
	Description string `form:"description" validate:"max=300"`
	// 最小起送金额,使用整数存储,默认存储到分
	MinimumDeliveryAmount uint `form:"minimum_delivery_amount"`
}

func (r *restaurantRequest) GetRestaurantModel() *models.Restaurant {
	res := models.Restaurant{
		RestaurantName:        r.RestaurantName,
		Address:               r.Address,
		Description:           r.Description,
		MinimumDeliveryAmount: r.MinimumDeliveryAmount,
	}
	return &res
}

// 应保证经过handler判定对应restaurant存在；
func UpdateRestaurant(c *gin.Context) {
	restaurant_id := c.GetUint("restaurant_id")
	var data restaurantRequest
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

	r := data.GetRestaurantModel()
	err = models.UpdateRestaurant(restaurant_id, r)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	app.ResponseSuccess(c)
}

// 该函数判定店铺名是否重复，如果重复则不允许创建
// 创建时后端只校验最大长度；前端负责校验其他
func CreateRestaurant(c *gin.Context) {
	merchant_id := handler.NewMerchInfoHanlder(c).GetID()

	var data restaurantRequest
	err := c.Bind(&data)
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return
	}
	// 注意添加外键
	err = app.ValidateStruct(data)
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return
	}
	r := data.GetRestaurantModel()
	r.MerchantID = merchant_id

	// 如果店铺存在那么就不应该再次创建
	exist, err := models.ExistRestaurant(r.RestaurantName)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	if exist {
		app.Response(c, http.StatusOK, ecode.ERROR_RESTAURANT_EXIST, nil)
		return
	}

	err = models.CreateRestaurant(r)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	app.ResponseSuccess(c)
}

func DeleteRestaurant(c *gin.Context) {
	// 因为已经提前设置在上下文中
	restaurant_id := c.GetUint("restaurant_id")

	err := models.DeleteRestaurant(restaurant_id)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccess(c)
}
