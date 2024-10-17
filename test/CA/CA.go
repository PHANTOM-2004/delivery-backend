package CA

import (
	"delivery-backend/internal/setting"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		log.Println("got [/]")
		c.JSON(http.StatusOK, gin.H{
			"message": "[/]test",
		})
	})
	r.GET("/test", func(c *gin.Context) {
		log.Println("got [/test]")
		c.JSON(http.StatusOK, gin.H{
			"message": "[/test]test",
		})
	})
	return r
}

// NOTE: 这里仅用于本地调试时开启，检测是否可以正常向localhost发起https请求
// 上线时会关闭
func LaunchServer() {
	router := InitRouter()
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.TestSetting.HTTPSPort),
		Handler:        router,
		MaxHeaderBytes: 1 << 20,
	}
	// s.ListenAndServe()
	certFile := "localhost-cert.pem"
	keyFile := "localhost-key.pem"
	log.Infof("test server for HTTPS test launching")
	log.Infof("listening port[%d]", setting.TestSetting.HTTPSPort)

	err := s.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}
}
