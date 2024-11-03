package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 传入page_id, 作为url传送
// 由管理员获取
func GetMerchantApplication(c *gin.Context) {
	page := c.Param("page")
	page_cnt, err := strconv.Atoi(page)
	if err != nil {
		app.Response(c, http.StatusOK, ecode.INVALID_PARAMS, nil)
		return
	}

	res, err := models.GetMerchantApplication(page_cnt)
	if err != nil {
		app.Response(c, http.StatusInternalServerError, ecode.ERROR, nil)
		return
	}

	data := map[string]any{
		"applications": res,
	}

	app.Response(c, http.StatusOK, ecode.SUCCESS, data)
}
