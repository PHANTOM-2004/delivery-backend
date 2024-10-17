package main

import (
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

func main() {
	// NOTE: 这里仅用于本地调试时开启，检测是否可以正常向localhost发起https请求
	// 上线时会关闭
	if true {
		router := InitRouter()
		s := &http.Server{
			Addr:           ":9000",
			Handler:        router,
			MaxHeaderBytes: 1 << 20,
		}
		// s.ListenAndServe()
		certFile := "localhost-cert.pem"
		keyFile := "localhost-key.pem"
		err := s.ListenAndServeTLS(certFile, keyFile)
		if err != nil {
			log.Fatal(err)
		}
	}
}
