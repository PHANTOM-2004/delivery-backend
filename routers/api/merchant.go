package api

import "github.com/gin-gonic/gin"

func MerchantGetAuth(c *gin.Context) {
	// 通过refresh_token, 获得access_token
	//
	// account := c.GetString("jwt_account")
}
