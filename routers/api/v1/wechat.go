package v1

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/setting"
	"delivery-backend/middleware/wechat"
	"delivery-backend/models"
	wechat_service "delivery-backend/service/wechat"
	"encoding/json"
	"net/http"
	"strconv"

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
	user, created, err := models.GetOrCreateWechatUser(openid)
	// 设置session
	session_id := uuid.NewString()
	session := wechat.NewWXSession(session_id)
	session.SetExpire()
	err = session.SetInfo(openid, user.ID, user.Role)
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	log.Tracef("session created: [%s]", session_id)

	if !created {
		// 应当返回一些具体信息
		app.ResponseSuccessWithData(c, map[string]any{
			"info":       user,
			"session_id": session_id,
		})
		return
	}

	log.Debug("created new user")
	// 接下来返回session
	app.ResponseSuccessWithData(c, map[string]any{
		"session_id": session_id,
	})
}

type UserInfoRequest struct {
	PhoneNumber     string `json:"phone_number" validate:"required,e164"`
	ProfileImageURL string `json:"profile_image_url" validate:"max=200"`
	NickName        string `json:"nickname" validate:"max=50"`
}

func (u *UserInfoRequest) GetModel() *models.WechatUser {
	return &models.WechatUser{
		PhoneNumber:     u.PhoneNumber,
		ProfileImageURL: u.ProfileImageURL,
		NickName:        u.NickName,
	}
}

func WXUploadUserInfo(c *gin.Context) {
	req := UserInfoRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.ResponseInvalidParams(c)
		log.Debug(err)
		return
	}
	err = app.ValidateStruct(&req)
	if err != nil {
		app.ResponseInvalidParams(c)
		log.Debug(err)
		return
	}

	session := wechat.DefaultSession(c)
	w, err := session.GetInfo()
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	err = models.UpdateWechatUser(w.ID, req.GetModel())
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccess(c)
}

func WXGetRestaurants(c *gin.Context) {
	res, err := models.GetAllRestaurants()
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccessWithData(c, map[string]any{
		"restaurants": res,
	})
}

func WXGetTopDishes(c *gin.Context) {
	restaurant_id, err := strconv.Atoi(c.Param("restaurant_id"))
	if err != nil {
		app.ResponseInvalidParams(c)
		log.Debug(err)
		return
	}

	res, err := models.GetTopDishes(uint(restaurant_id))
	if err != nil {
		app.ResponseInternalError(c, err)
		return
	}
	app.ResponseSuccessWithData(c, map[string]any{
		"dishes": res,
	})
}
