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

// 传入page_id, 作为url传送
// 由管理员获取
func GetMerchantApplication(c *gin.Context) {
	page := c.Param("page")
	page_cnt, err := strconv.Atoi(page)
	if err != nil {
		log.Debug(err)
		app.Response(c, http.StatusOK, ecode.INVALID_PARAMS, nil)
		return
	}

	res, err := models.GetMerchantApplications(page_cnt)
	if err != nil {
		app.Response(c, http.StatusInternalServerError, ecode.ERROR, nil)
		return
	}

	data := map[string]any{
		"applications": res,
	}

	app.Response(c, http.StatusOK, ecode.SUCCESS, data)
}

func ApproveMerchantApplication(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Debug(err)
		app.Response(c, http.StatusOK, ecode.INVALID_PARAMS, nil)
		return
	}

	// status从 未审核变为通过审核，或者从不通过审核变为通过审核
	// 如果原本已经通过审核，则该请求不造成任何后果
	//
	//
	// 需要判定先前状态, 这种申请必然不频繁，所以直接从数据库查询即可
	// 至于先前状态，实际上只需要找商家账号是否存在即可
	m, err := models.GetRelatedMerchant(id)
	if err != nil {
		log.Warn(err)
		app.Response(c, http.StatusInternalServerError, ecode.ERROR, nil)
		return
	}

	if m != nil {
		// 说明商家有账号，所以需要检查是否解冻这个账号
		// TODO:通过redis黑名单+db冻结
	}

	// 更新状态
	err = models.ApproveApplication(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Warnf("application form id[%d] not found", id)
		app.Response(c, http.StatusOK, ecode.ERROR_MERCHANT_APPLICATION_NOT_FOUND, nil)
		return
	} else if err != nil {
		log.Warn(err)
		app.Response(c, http.StatusInternalServerError, ecode.ERROR, nil)
		return
	}

	app.Response(c, http.StatusOK, ecode.SUCCESS, nil)

	log.Debugf("application form id[d] approved", id)
}

func DisapproveMerchantApplication(c *gin.Context) {
}
