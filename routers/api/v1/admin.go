package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/models"
	"delivery-backend/service/merchant_service"
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
		app.ResponseInvalidParams(c)
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
	id, err := strconv.Atoi(c.Param("application_id"))
	if err != nil {
		log.Warn(err)
		app.ResponseInvalidParams(c)
		return
	}

	// status从 未审核变为通过审核，或者从不通过审核变为通过审核
	// 如果原本已经通过审核，则该请求不造成任何后果
	//
	//
	// 需要判定先前状态, 这种申请必然不频繁，所以直接从数据库查询即可
	// 至于先前状态，实际上只需要找商家账号是否存在即可
	//
	// 1. 查找相关联的账号
	m, err := models.GetRelatedMerchant(uint(id))
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	// 2. 存在商家账号
	if m.ID != 0 {
		// 说明商家有账号，所以需要检查是否解冻这个账号
		log.Debugf("related account[%s] found", m.Account)
		err = merchant_service.EnableMerchant(m.ID)
		if err != nil {
			app.ResponseInternalError(c, err)
			return
		}
		log.Debugf("merchant account[%s] enabled", m.Account)
	}

	// 3.更新申请表的状态
	err = models.ApproveApplication(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Warnf("application form id[%d] not found", id)
		app.Response(c, http.StatusOK,
			ecode.ERROR_MERCHANT_APPLICATION_NOT_FOUND, nil)
		return
	} else if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	// 商家已经存在，那么不需要考虑创建账号
	if m.ID != 0 {
		app.ResponseSuccess(c)
		log.Debugf("application form id[%d] approved", id)
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

func DisapproveMerchantApplication(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("application_id"))
	if err != nil {
		log.Debug(err)
		app.ResponseInvalidParams(c)
		return
	}
	m, err := models.GetRelatedMerchant(uint(id))
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	if m.ID == 0 {
		// 没有关联的商家账号，那么那么无需冻结和解冻操作
		log.Debug("No related merchant account")
		app.ResponseSuccess(c)
		return
	}

	// 有关联的商家账号，需要冻结
	// NOTE:正常来说，关联的商家此时必然是非冻结状态，否则就是多次disapprove
	err = merchant_service.DisableMerchant(m.ID)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	log.Debugf("Disable merchant account: %s", m.Account)
	app.ResponseSuccess(c)
}
