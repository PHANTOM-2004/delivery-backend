package wechat

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/internal/gredis"
	"delivery-backend/internal/setting"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const ginDefaultWXSessionKey = "CYT_WX_Session"

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

// 设置key的过期时间
func (wxs *WXSession) SetExpire() error {
	expires := time.Duration(setting.WechatSetting.SessionAge) * time.Second
	err := gredis.Expire(wxs.session_id, expires)
	return err
}

func DefaultSession(c *gin.Context) *WXSession {
	res, exist := c.Get(ginDefaultWXSessionKey)
	if !exist {
		log.Panic("no wx session, it is the fault of coder")
	}
	return res.(*WXSession)
}

type WXSessionCartStore struct {
	DishID   uint `json:"dish_id"`
	FlavorID uint `json:"flavor_id"`
	Cnt      int  `json:"count"`
}



func (wxs *WXSession) getCartKey(restaurant_id uint) string {
	return "cart_" + strconv.Itoa(int(restaurant_id))
}

func (wxs *WXSession) UpdateCart(restaurant_id uint, store []WXSessionCartStore) error {
	field := wxs.getCartKey(restaurant_id)
	s, err := json.Marshal(&store)
	if err != nil {
		log.Panic(err)
	}
	err = gredis.HSet(wxs.session_id, []string{field, string(s)})
	return err
}

func (wxs *WXSession) GetCart(restaurant_id uint) ([]WXSessionCartStore, error) {
	field := wxs.getCartKey(restaurant_id)
	res, err := gredis.HGet(wxs.session_id, field)
	if err != nil {
		return nil, err
	}

	cart := []WXSessionCartStore{}
	err = json.Unmarshal([]byte(res), &cart)
	if err != nil {
		log.Panic(err)
	}
	return cart, err
}

func (wxs *WXSession) SetInfo(openid string, id uint, role uint8) error {
	pairs := []string{
		"openid", openid,
		"id", strconv.Itoa(int(id)),
		"role", strconv.Itoa(int(role)),
	}
	err := gredis.HSet(wxs.session_id, pairs)
	return err
}

func (wxs *WXSession) GetInfo() (*WXSessionInfoStore, error) {
	var err error
	var openid, role, id string
	openid, err1 := gredis.HGet(wxs.session_id, "openid")
	id, err2 := gredis.HGet(wxs.session_id, "id")
	role, err3 := gredis.HGet(wxs.session_id, "role")
	if err1 != nil || err2 != nil || err3 != nil {
		return nil, fmt.Errorf("[%v]+[%v]+[%v]", err1, err2, err3)
	}

	r_id, err := strconv.Atoi(id)
	if err != nil {
		log.Panic("must not fail")
	}
	r_role, err := strconv.Atoi(role)
	if err != nil {
		log.Panic("must not fail")
	}

	res := WXSessionInfoStore{
		OpenID: openid,
		Role:   uint8(r_role),
		ID:     uint(r_id),
	}
	return &res, err
}

// 从中提取出session_id
func WXsession() gin.HandlerFunc {
	return func(c *gin.Context) {
		// handle 微信发送的session id 请求
		session_id := c.GetHeader("Authorization")
		if session_id == "" {
			log.Debug("没有在请求头中提供session_id")
			app.ResponseInvalidParams(c)
			return
		}
		log.Tracef("get session_id from header [%s]", session_id)
		// 检查session是否过期，如果不存在说明已经过期
		exist, err := gredis.Exists(session_id)
		if err != nil {
			app.ResponseInternalError(c, err)
			return
		}
		if !exist {
			// 说明session已经过期
			app.Response(c, http.StatusOK, ecode.ERROR_WX_SESSION_EXPIRE, nil)
			c.Abort()
			return
		}
		// 如果session没有过期， 设置session
		c.Set(ginDefaultWXSessionKey, NewWXSession(session_id))
		log.Trace("Set Session in Context")
		c.Next()
	}
}
