package jwt

import (
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type (
	TokenInBlacklistFunc = func(string) bool
	TokenAuthFunc        = func(string) (uint, string, ecode.Ecode)
)

// 在身份认证中，set或者get对应的id和账户
type JwtInfo struct {
	c *gin.Context
}

func NewJwtInfo(c *gin.Context) *JwtInfo {
	return &JwtInfo{c}
}

func (j *JwtInfo) GetID() uint {
	return j.c.GetUint("jwt_id")
}

func (j *JwtInfo) GetAccount() uint {
	return j.c.GetUint("jwt_account")
}

func (j *JwtInfo) SetID(id uint) {
	j.c.Set("jwt_id", id)
}

func (j *JwtInfo) SetAccount(account string) {
	j.c.Set("jwt_account", account)
}

func JWTRK(b TokenInBlacklistFunc, a TokenAuthFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// cookie中读取refresh_token
		refresh_token, err := c.Cookie("refresh_token")
		if errors.Is(err, http.ErrNoCookie) {
			defer log.Debug("fail: no refresh_token")
			defer c.Abort()
			app.Response(c, http.StatusUnauthorized, ecode.ERROR_AUTH_NO_REFRESH_TOKEN, nil)
			return
		}

		log.Debug("pass: refresh_token exists")

		// 校验refresh token
		id, account, code := a(refresh_token)
		if code != ecode.SUCCESS {
			defer c.Abort()
			app.Response(c, http.StatusUnauthorized, code, nil)
			return
		}

		log.Debug("pass: refresh_token validation")

		// 检查refresh token是否在黑名单中
		in_blacklist := b(refresh_token)
		if in_blacklist {
			defer log.Debug("fail: refresh_token in blacklist")
			defer c.Abort()
			app.Response(c, http.StatusUnauthorized, ecode.ERROR_AUTH_REFRESH_TOKEN_EXPIRED, nil)
			return
		}

		log.Debug("pass: refresh_token not in blacklist")

		j := NewJwtInfo(c)
		j.SetAccount(account)
		j.SetID(id)

		c.Next()
	}
}

func JWTAK(b TokenInBlacklistFunc, a TokenAuthFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		access_token, err := c.Cookie("access_token")
		if errors.Is(err, http.ErrNoCookie) {
			app.Response(c, http.StatusUnauthorized, ecode.ERROR_AUTH_NO_ACCESS_TOKEN, nil)
			c.Abort()
			log.Debug("fail: no access_token provided")
			return
		}

		log.Debug("pass: received jwt access_token")

		id, account, code := a(access_token)
		if code != ecode.SUCCESS {
			app.Response(c, http.StatusUnauthorized, code, nil)
			c.Abort()
			return
		}

		log.Debug("pass: access_token validation")

		// 检查是否在黑名单中
		in_blacklist := b(access_token)
		if in_blacklist {
			app.Response(c, http.StatusUnauthorized, ecode.ERROR_AUTH_ACCESS_TOKEN_EXPIRED, nil)
			c.Abort()
			log.Debug("fail: access_token in blacklist")
			return
		}

		c.Set("jwt_account", account)
		c.Set("jwt_id", id)

		log.Debug("pass: access_token not in blacklist")

		c.Next()
	}
}
