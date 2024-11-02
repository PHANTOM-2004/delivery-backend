package service

import (
	"delivery-backend/internal/ecode"
	"delivery-backend/internal/setting"
	jwt_token "delivery-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func AuthAdminAccessToken(access_token string) (string, ecode.Ecode) {
	return jwt_token.AuthAccessToken("admin", access_token)
}

func AuthAdminRefreshToken(refresh_token string) (string, ecode.Ecode) {
	return jwt_token.AuthRefreshToken("admin", refresh_token)
}

// 将token加入黑名单，同时删除cookie中token
func DisbleAdminTokens(c *gin.Context) {
	jwt_token.DisableTokens("ADMIN_TK_", c)
}

// 判断token是否在redis黑名单中
// 当token有效返回true，无效返回false
func ValidateAdminToken(admin_token string) bool {
	return jwt_token.ValidateToken("ADMIN_TK_", admin_token)
}

func SetRefreshToken(c *gin.Context, refresh_token string) {
	jwt_token.SetRefreshToken(c, refresh_token, setting.AppSetting.AdminRKAge)
}

func SetAccessToken(c *gin.Context, access_token string) {
	jwt_token.SetRefreshToken(c, access_token, setting.AppSetting.AdminAKAge)
}
