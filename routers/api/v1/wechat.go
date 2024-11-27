package v1

import (
	"crypto/tls"
	"delivery-backend/internal/app"
	"delivery-backend/internal/setting"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var WXClient = &http.Client{
	Timeout: 10 * time.Second, // 设置超时时间
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			// 跳过证书验证（仅限测试，生产环境中请使用有效证书）
		},
	},
}

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
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}

	resp, err := WXClient.Do(req)
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
	app.ResponseSuccess(c)
}
