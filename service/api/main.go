package main

import (
	"delivery-backend/service/api/middleware"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func InitRouter() *gin.Engine {
	defer log.Info("app router initialized")

	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	// TODO:
	// r.MaxMultipartMemory = int64(setting.AppSetting.MaxImageSize << 20) // MiB

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("[GIN-route] %6s %v -- [%v] (%v handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	gin.DebugPrintFunc = func(format string, values ...any) {
		log.Debug("[GIN-debug] "+format, values)
	}

	apiv1 := r.Group("/api/v1")
	initMerchantRouter(apiv1)

	return r
}

func LaunchServer() {
	// launch server
	r := InitRouter()
	s := &http.Server{
		// TODO:
		Addr:           fmt.Sprintf(":%d", 8000),
		Handler:        r,
		MaxHeaderBytes: 1 << 20,
	}
	log.Infof("listening port[%d]", 8000)
	s.ListenAndServe()
}

func main() {
	LaunchServer()
}
