package wechat

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/internal/gredis"
	"delivery-backend/internal/setting"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const defaultWXSessionKey = "CYT_WX_Session"

type WXSessionInfoStore struct {
	OpenID string `json:"openid"`
	ID     uint   `json:"id"`
	Role   uint8  `json:"role"`
}

type WXSession struct {
	session_id string
}

func NewWXSession(session_id string) *WXSession {
	res := WXSession{
		session_id: session_id,
	}
	return &res
}

func DefaultSession(c *gin.Context) *WXSession {
	res, exist := c.Get(defaultWXSessionKey)
	if !exist {
		log.Fatal("no wx session, fatal error")
	}
	return res.(*WXSession)
}

func (wxs *WXSession) SetInfo(openid string, id uint, role uint8) error {
	s := WXSessionInfoStore{
		OpenID: openid,
		ID:     id,
		Role:   role,
	}
	expires := time.Duration(setting.WechatSetting.SessionAge) * time.Second
	err := gredis.Set(wxs.session_id, s, expires)
	return err
}

func (wxs *WXSession) GetInfo(openid string, id uint) (*WXSessionInfoStore, error) {
	s, err := gredis.Get(wxs.session_id)
	if err != nil {
		return nil, err
	}
	res := WXSessionInfoStore{}
	err = json.Unmarshal(s, &res)
	return &res, err
}

// 从中提取出session_id
func WXsession() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO:
		// handle 微信发送的session id 请求
		session_id := c.Query("session_id")
		if session_id == "" {
			log.Debug("没有提供session_id")
			app.ResponseInvalidParams(c)
			return
		}
		// 检查session是否过期，如果不存在说明已经过期
		exist, err := gredis.Exists(session_id)
		if err != nil {
			app.ResponseInternalError(c, err)
			return
		}
		if !exist {
			// 说明session已经过期
			app.Response(c, http.StatusOK, ecode.ERROR_WX_SESSION_EXPIRE, nil)
			return
		}
		// 如果session没有过期， 设置session
		c.Set(defaultWXSessionKey, NewWXSession(session_id))
		c.Next()
	}
}
