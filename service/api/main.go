package main

import (
	"delivery-backend/service/api/biz/dal"
	"delivery-backend/service/api/conf"
	"delivery-backend/service/api/middleware"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
	"github.com/gin-gonic/gin"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap/zapcore"
)

func InitRouter() *gin.Engine {
	defer klog.Info("app router initialized")

	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.IPRateLimiter())
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong!")
	})

	// TODO:
	// r.MaxMultipartMemory = int64(setting.AppSetting.MaxImageSize << 20) // MiB

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		klog.Infof("[GIN-route] %6s %v -- [%v] (%v handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	gin.DebugPrintFunc = func(format string, values ...any) {
		klog.Debug("[GIN-debug] "+format, values)
	}

	apiv1 := r.Group("/api/v1")
	initMerchantRouter(apiv1)
	initAdminRouter(apiv1)

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
	klog.Infof("listening port[%d]", 8000)
	s.ListenAndServe()
}

func InitKlog() {
	// klog
	logger := kitexlogrus.NewLogger()
	klog.SetLogger(logger)
	klog.SetLevel(conf.LogLevel())
	asyncWriter := &zapcore.BufferedWriteSyncer{
		WS: zapcore.AddSync(&lumberjack.Logger{
			Filename:   conf.GetConf().Kitex.LogFileName,
			MaxSize:    conf.GetConf().Kitex.LogMaxSize,
			MaxBackups: conf.GetConf().Kitex.LogMaxBackups,
			MaxAge:     conf.GetConf().Kitex.LogMaxAge,
		}),
		FlushInterval: time.Second * 3,
	}
	w := io.MultiWriter(os.Stdout, asyncWriter)
	klog.SetOutput(w)

	server.RegisterShutdownHook(func() {
		asyncWriter.Sync()
	})
	klog.Info("klog init done")
	klog.Debug("klog init done")
	slog.Debug("klog init done")
	fmt.Println("klog")
}

func main() {
	dal.Init()
	InitKlog()
	LaunchServer()
}
