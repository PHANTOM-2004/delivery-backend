package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/internal/setting"
	"delivery-backend/models"
	"delivery-backend/service/customer_service"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func getLicenseFileName() string {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Warn(err)
	}
	path := setting.AppSetting.LicenseStorePath + "/merchant-license-" + id.String()
	return path
}

func checkLicenseType(name string) bool {
	ext := filepath.Ext(name)
	allows := setting.AppSetting.LicenseAllowExts
	for i := range allows {
		if ext == allows[i] {
			return true
		}
	}
	return false
}

// 由顾客/骑手发起申请，提出商务合作请求
// https://gin-gonic.com/docs/examples/upload-file/single-file/
func MerchantApply(c *gin.Context) {
	// TODO: 顾客身份校验

	file, err := c.FormFile("license")
	if err != nil || !checkLicenseType(file.Filename) {
		log.Debug(err)
		app.Response(c, http.StatusOK, ecode.INVALID_PARAMS, nil)
		return
	}

	a := customer_service.Application{
		Description: c.PostForm("description"),
		Email:       c.PostForm("email"),
		PhoneNumber: c.PostForm("phone_number"),
	}

	err = app.ValidateStruct(a)
	if err != nil {
		log.Debug(err)
		app.Response(c, http.StatusOK, ecode.INVALID_PARAMS, nil)
		return
	}

	log.Debug("uploaded: ", file.Filename)

	// 保存证书，并且使用id重命名
	path := getLicenseFileName()
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		log.Warn(err)
		app.Response(c, http.StatusInternalServerError, ecode.ERROR, nil)
		return
	}

	// 插入数据库

	data := models.MerchantApplication{
		Description: a.Description,
		License:     path,
		Email:       a.Email,
		PhoneNumber: a.PhoneNumber,
	}

	err = models.CreateMerchantApplication(&data)
	if err != nil {
		log.Warn(err)
		app.Response(c, http.StatusInternalServerError, ecode.ERROR, nil)
		return
	}

	// 成功保存
	log.Debugf("license[%s] saved to %s", file.Filename, path)

	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)
}
