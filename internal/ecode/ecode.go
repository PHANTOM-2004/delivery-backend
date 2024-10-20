package ecode

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Ecode = int

const (
	SUCCESS        Ecode = http.StatusOK
	ERROR          Ecode = http.StatusInternalServerError
	INVALID_PARAMS Ecode = 1001

	// admin related error
	ERROR_ADMIN_NON_EXIST     Ecode = 1002
	ERROR_ADMIN_INCORRECT_PWD Ecode = 1003

	ERROR_AUTH_CHECK_TOKEN_FAIL    Ecode = 20001
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT Ecode = 20002
	ERROR_AUTH_TOKEN               Ecode = 20003
	ERROR_AUTH                     Ecode = 20004
)

func StatusText(e Ecode) (res string) {
	switch e {
	case SUCCESS:
		res = "ok"
	case ERROR:
		res = "fail"
	case ERROR_ADMIN_NON_EXIST:
		res = "管理员账号不存在"
	case ERROR_ADMIN_INCORRECT_PWD:
		res = "管理员密码输入错误"
	case INVALID_PARAMS:
		res = "请求参数错误"
	case ERROR_AUTH_CHECK_TOKEN_FAIL:
		res = "Token鉴权失败"
	case ERROR_AUTH_CHECK_TOKEN_TIMEOUT:
		res = "Token已超时"
	case ERROR_AUTH_TOKEN:
		res = "Token生成失败"
	case ERROR_AUTH:
		res = "Token错误"
	default:
		log.Fatalf("未知错误码[%d]", e)
	}
	return
}
