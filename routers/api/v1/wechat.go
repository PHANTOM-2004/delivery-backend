package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/setting"
	"delivery-backend/middleware/wechat"
	"delivery-backend/models"
	wechat_service "delivery-backend/service/wechat"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	// 接下来如果用户不存在，那么就创建用户；否则就返回正确用户
	openid := wxserverResp["openid"].(string)
	user, err := models.GetOrCreateWechatUser(openid)
	// 设置session
	session_id := uuid.NewString()
	session := wechat.NewWXSession(session_id)
	err = session.SetInfo(openid, user.ID, user.Role)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	log.Tracef("session created: [%s]", session_id)

	// 接下来返回session
	app.ResponseSuccessWithData(c, map[string]any{
		"session_id": session_id,
	})
}
