package admin_service

import (
	"delivery-backend/internal/ecode"
	"delivery-backend/internal/setting"
	jwt_token "delivery-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func AuthAccessToken(access_token string) (uint, string, ecode.Ecode) {
	c, ecode := jwt_token.AuthAccessToken("admin", access_token)
	if c != nil {
		return c.ID, c.Account, ecode
	}
	return 0, "", ecode
}

func AuthRefreshToken(refresh_token string) (uint, string, ecode.Ecode) {
	c, ecode := jwt_token.AuthRefreshToken("admin", refresh_token)
	if c != nil {
		return c.ID, c.Account, ecode
	}
	return 0, "", ecode
}

// 将token加入黑名单，同时删除cookie中token
func DisbleTokens(c *gin.Context) {
	jwt_token.DisableTokens("ADMIN_TK_", c)
}

// 判断token是否在redis黑名单中
// 当token在黑名单中返回true，否则返回false
func TokenInBlacklist(admin_token string) bool {
	return jwt_token.TokenInBlacklist("ADMIN_TK_", admin_token)
}

func SetRefreshToken(c *gin.Context, refresh_token string) {
	jwt_token.SetRefreshToken(c, refresh_token, setting.AppSetting.AdminRKAge)
}

func SetAccessToken(c *gin.Context, access_token string) {
	jwt_token.SetAccessToken(c, access_token, setting.AppSetting.AdminAKAge)
}

func GetAccessToken(id uint, account string) string {
	return jwt_token.GetAccessToken("admin", id, account, setting.AppSetting.AdminAKAge)
}

func GetRefreshToken(id uint, account string) string {
	return jwt_token.GetRefreshToken("admin", id, account, setting.AppSetting.AdminRKAge)
}
