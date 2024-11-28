package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/internal/setting"
	"delivery-backend/models"
	"delivery-backend/service/customer_service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// 由顾客/骑手发起申请，提出商务合作请求
// https://gin-gonic.com/docs/examples/upload-file/single-file/
func MerchantApply(c *gin.Context) {
	// TODO: 顾客身份校验

	file, err := c.FormFile("license")
	if err != nil {
		app.ResponseInvalidParams(c)
		return
	}
	ext, v := setting.AppSetting.CheckLicenseExt(file.Filename)
	if !v {
		log.Debugf("wrong ext[%s]", ext)
		app.ResponseInvalidParams(c)
		return
	}

	a := customer_service.Application{
		Description: c.PostForm("description"),
		Email:       c.PostForm("email"),
		PhoneNumber: c.PostForm("phone_number"),
		Name:        c.PostForm("name"),
	}

	err = app.ValidateStruct(a)
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return
	}

	log.Debug("uploadint: ", file.Filename)
	// 保存证书，并且使用id重命名,更改回ext
	name := setting.AppSetting.GenLicenseName() + ext
	dst := setting.AppSetting.GetLicenseStorePath(name)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	log.Trace("saved to:", dst)

	// 插入数据库

	data := models.MerchantApplication{
		Description: a.Description,
		License:     name,
		Email:       a.Email,
		PhoneNumber: a.PhoneNumber,
		Name:        a.Name,
	}

	err = models.CreateMerchantApplication(&data)
	if err != nil {
		log.Warn(err)
		app.Response(c, http.StatusInternalServerError, ecode.ERROR, nil)
		return
	}

	// 成功保存
	log.Debugf("license[%s] saved to %s", file.Filename, dst)

	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
}

func GetRestaurantCategoryDish(c *gin.Context) {
	var err error
	restaurant_id, err := strconv.Atoi(c.Param("restaurant_id"))
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return
	}

	res, err := models.GetCategoryDishFlavor(uint(restaurant_id))
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccessWithData(c, map[string]any{"categories": res})
}
