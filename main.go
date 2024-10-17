package main

import (
	"delivery-backend/internal/setting"
	"delivery-backend/test/CA"
)

func main() {
	setting.Setup()

	if setting.TestSetting.CATest {
    //NOTE: 请使用该模块测试本地开发环境是否可以正确访问
    //https://localhost:xxxx
		go CA.LaunchServer()
		// block main function
		select {}
	}
}
