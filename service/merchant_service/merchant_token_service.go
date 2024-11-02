package merchant_service

import (
	"delivery-backend/internal/ecode"
	"delivery-backend/internal/setting"
	jwt_token "delivery-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func AuthAccessToken(access_token string) (string, ecode.Ecode) {
	return jwt_token.AuthAccessToken("merchant", access_token)
}

func AuthRefreshToken(refresh_token string) (string, ecode.Ecode) {
	return jwt_token.AuthRefreshToken("merchant", refresh_token)
}

// 将token加入黑名单，同时删除cookie中token
func DisbleTokens(c *gin.Context) {
	jwt_token.DisableTokens("MERCH_TK_", c)
}

// 判断token是否在redis黑名单中
// 当token有效返回true，无效返回false
func ValidateToken(admin_token string) bool {
	return jwt_token.ValidateToken("MERCH_TK_", admin_token)
}

func SetRefreshToken(c *gin.Context, refresh_token string) {
	jwt_token.SetRefreshToken(c, refresh_token, setting.AppSetting.MerchantRKAge)
}

func SetAccessToken(c *gin.Context, access_token string) {
	jwt_token.SetRefreshToken(c, access_token, setting.AppSetting.MerchantAKAge)
}

func GetAccessToken(account string) string {
	return jwt_token.GetAccessToken("merchant", account, setting.AppSetting.MerchantAKAge)
}

func GetRefreshToken(account string) string {
	return jwt_token.GetRefreshToken("merchant", account, setting.AppSetting.MerchantRKAge)
}
