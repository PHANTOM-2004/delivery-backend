package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/setting"
	wechat_service "delivery-backend/service/wechat"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type WXLoginRequest struct {
	Code string `json:"code"`
}

func WXLogin(c *gin.Context) {
	var wxrequest WXLoginRequest
	err := c.ShouldBindJSON(&wxrequest)
	if err != nil {
		app.ResponseInvalidParams(c)
		return
	}
	log.Trace(wxrequest)
	code := wxrequest.Code
	url := setting.WechatSetting.GetCode2SessionURL(code)
	log.Tracef("sending code[%s] to [%s]", code, url)

	// TODO: handle error

	// send request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	resp, err := wechat_service.WXClient.Do(req)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	wxserverResp := map[string]any{}
	err = json.NewDecoder(resp.Body).Decode(&wxserverResp)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	log.Trace("wxserver code2sesion response: ", wxserverResp)
	// 接下来返回session
	app.ResponseSuccess(c)
}
