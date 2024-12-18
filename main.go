package main

import (
	"delivery-backend/internal/gredis"
	"delivery-backend/internal/setting"
	"delivery-backend/models"
	"delivery-backend/routers"
	"delivery-backend/service/consumer"
	"delivery-backend/service/email"
	wechat_service "delivery-backend/service/wechat"
	"delivery-backend/test/CA"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Setup() {
	setting.Setup()
	models.SetUp()
	gredis.Setup()
	// set server mode
	gin.SetMode(setting.ServerSetting.RunMode)

	// 初始化consumer email
	email.Setup()
	// 启动所有consumer
	consumer.Setup()

	// 初始化wechat access token服务
	wechat_service.Setup()
}

// shutdown服务
func Shutdown() {
	defer consumer.ShutDown()
}

func LaunchServer() {
	// launch server
	r := routers.InitRouter()
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.ServerSetting.HTTPPort),
		Handler:        r,
		MaxHeaderBytes: 1 << 20,
	}
	log.Infof("listening port[%d]", setting.ServerSetting.HTTPPort)
	s.ListenAndServe()
	// certFile := setting.ServerSetting.SSLCertPath
	// keyFile := setting.ServerSetting.SSLKeyPath
	// err := s.ListenAndServeTLS(certFile, keyFile)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func main() {
	Setup()
	defer Shutdown()

	if setting.TestSetting.CATest && gin.Mode() == gin.DebugMode {
		// NOTE: 请使用该模块测试本地开发环境是否可以正确访问
		// https://localhost:xxxx
		go CA.LaunchServer()
	}

	LaunchServer()
}
