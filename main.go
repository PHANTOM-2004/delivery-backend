package main

import (
	"delivery-backend/internal/server"
	"delivery-backend/internal/setting"
	"delivery-backend/test/CA"

	"github.com/gin-gonic/gin"
)

func main() {
	server.Setup()

	if setting.TestSetting.CATest && gin.Mode() == gin.DebugMode {
		// NOTE: 请使用该模块测试本地开发环境是否可以正确访问
		// https://localhost:xxxx
		go CA.LaunchServer()
	}

	server.LaunchServer()
}
