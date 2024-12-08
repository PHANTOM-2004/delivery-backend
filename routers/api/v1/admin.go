package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/internal/setting"
	"delivery-backend/models"
	"delivery-backend/pkg/utils"
	"delivery-backend/service/merchant_service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// 传入page_id, 作为url传送
// 由管理员获取
func GetMerchantApplication(c *gin.Context) {
	page := c.Param("page")
	page_cnt, err := strconv.Atoi(page)
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return
	}

	res, err := models.GetMerchantApplications(page_cnt)
	if err != nil {
		app.Response(c, http.StatusBadRequest, ecode.ERROR, nil)
		return
	}

	data := map[string]any{
		"applications": res,
	}

	app.Response(c, http.StatusOK, ecode.SUCCESS, data)
}

// approve之前必然是拒绝或者没有看的状态, 也就是没有对应的账号, 所以approve的时候必然创建.
func ApproveMerchantApplication(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("application_id"))
	if err != nil {
		log.Warn(err)
		app.ResponseInvalidParams(c)
		return
	}

	// 3.更新申请表的状态
	succ, err := models.ApproveApplication(id)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	if !succ {
		app.Response(c, http.StatusOK, ecode.ERROR_ADMIN_INVALID_OPERATION, nil)
		return
	}

	// 4.考虑为商家注册账号
	err = merchant_service.CreateMerchantFromApplication(uint(id))
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	log.Debugf("application form id[%d] approved; account created", id)
	app.ResponseSuccess(c)
}

// 对于disapprove, 不再进行封禁操作, 也就是说disapprove后可以approve, approve后不能disapprove
func DisapproveMerchantApplication(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("application_id"))
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return
	}
	succ, err := models.DisapproveApplication(id)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	if !succ {
		app.Response(c, http.StatusBadRequest, ecode.ERROR_ADMIN_INVALID_OPERATION, nil)
		return
	}
	app.ResponseSuccess(c)
}

func DeleteMerchant(c *gin.Context) {
	account := c.PostForm("account")
	err, rows := models.DeleteMerchant(account)
	if err != nil {
		res := map[string]string{
			"error": err.Error(),
		}
		app.Response(c, http.StatusOK, ecode.ERROR, res)
		return
	}

	if rows <= 0 {
		res := map[string]string{
			"warn": "delete nothing",
		}
		app.Response(c, http.StatusOK, ecode.SUCCESS, res)
		return
	}

	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
}

func CreateMerchant(c *gin.Context) {
	if v := merchant_service.SignUpRequestValidate(c); !v {
		return
	}

	// create account
	account := c.PostForm("account")
	encrypted_password := utils.Encrypt(c.PostForm("password"), setting.AppSetting.Salt)
	merchant_name := c.PostForm("merchant_name")
	phone_numer := c.PostForm("phone_number")
	application_id, err := strconv.Atoi(c.PostForm("merchant_application_id"))
	if err != nil {
		app.ResponseInvalidParams(c)
		return
	}

	data := models.Merchant{
		MerchantName:          merchant_name,
		Account:               account,
		Password:              encrypted_password,
		PhoneNumber:           phone_numer,
		MerchantApplicationID: uint(application_id),
	}

	created, err := models.CreateMerchant(&data)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	if !created {
		app.Response(c, http.StatusOK, ecode.ERROR_MERCHANT_ACCOUNT_EXIST, nil)
		return
	}
	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
}

func GetMerchants(c *gin.Context) {
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		log.Warn(err)
		app.ResponseInvalidParams(c)
		return
	}
	merchants, err := models.GetMerchants(page)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	res := map[string]any{
		"merchants": merchants,
	}
	app.Response(c, http.StatusOK, ecode.SUCCESS, res)
}

func EnableMerchant(c *gin.Context) {
	merchant_id, err := strconv.ParseUint(c.Param("merchant_id"), 10, 0)
	if err != nil {
		log.Warn(err)
		app.ResponseInvalidParams(c)
		return
	}

	err = merchant_service.EnableMerchant(uint(merchant_id))
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	app.ResponseSuccess(c)
}

func DisableMerchant(c *gin.Context) {
	merchant_id, err := strconv.ParseUint(c.Param("merchant_id"), 10, 0)
	if err != nil {
		log.Warn(err)
		app.ResponseInvalidParams(c)
		return
	}

	err = merchant_service.DisableMerchant(uint(merchant_id))
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	app.ResponseSuccess(c)
}

func ApproveRiderApplication(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("application_id"))
	if err != nil {
		log.Warn(err)
		app.ResponseInvalidParams(c)
		return
	}
	succ, err := models.ApproveRider(uint(id))
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	if !succ {
		app.Response(c, http.StatusBadRequest, ecode.ERROR_ADMIN_INVALID_OPERATION, nil)
		return
	}
	app.ResponseSuccess(c)
}

// 对于disapprove, 不再进行封禁操作, 也就是说disapprove后可以approve, approve后不能disapprove
func DisapproveRiderApplication(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("application_id"))
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return
	}
	succ, err := models.DisapproveRider(uint(id))
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	if !succ {
		app.Response(c, http.StatusBadRequest, ecode.ERROR_ADMIN_INVALID_OPERATION, nil)
		return
	}
	app.ResponseSuccess(c)
}

func GetRiderApplications(c *gin.Context) {
	res, err := models.GetRiderApplications()
	if err != nil {
		app.ResponseInternalError(c, err)
		return

	}
	app.ResponseSuccessWithData(c, map[string]any{
		"applications": res,
	})
}

func GetRestaurantOrders(c *gin.Context) {
	restaurant_id, err := strconv.Atoi("restaurant_id")
	if err != nil {
		app.ResponseInvalidParams(c)
		return
	}
	orders, err := models.GetOrderByRestaurant(uint(restaurant_id))
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccessWithData(c, map[string]any{
		"orders": orders,
	})
}
