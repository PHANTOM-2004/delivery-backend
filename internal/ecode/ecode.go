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
	ERROR_ADMIN_NON_EXIST     Ecode = 10000
	ERROR_ADMIN_INCORRECT_PWD Ecode = 10001
	ERROR_ADMIN_ACCOUNT_EXIST Ecode = 10002
	ERROR_ADMIN_LOGOUT        Ecode = 10003
	ERROR_ADMIN_NOT_LOGIN     Ecode = 10004
	ERROR_ADMIN_ROLE          Ecode = 10005

	ERROR_AUTH_NO_TOKEN            Ecode = 20000
	ERROR_AUTH_CHECK_TOKEN_FAIL    Ecode = 20001
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT Ecode = 20002
	ERROR_AUTH_TOKEN_GENERATE      Ecode = 20003
	ERROR_AUTH                     Ecode = 20004

	ERROR_SUPER_AUTH          Ecode = 30000
	ERROR_SUPER_AUTH_NO_TOKEN Ecode = 30001
)

func StatusText(e Ecode) (res string) {
	switch e {
	case SUCCESS:
		res = "ok"
	case ERROR:
		res = "fail"
	case ERROR_ADMIN_NON_EXIST:
		res = "管理员账号不存在"
	case ERROR_ADMIN_ACCOUNT_EXIST:
		res = "管理员账号已注册"
	case ERROR_ADMIN_INCORRECT_PWD:
		res = "管理员密码输入错误"
	case ERROR_ADMIN_LOGOUT:
		res = "管理员非法登出请求/已登出"
	case ERROR_ADMIN_NOT_LOGIN:
		res = "管理员未登陆"
	case ERROR_ADMIN_ROLE:
		res = "管理员身份错误"
	case INVALID_PARAMS:
		res = "请求参数错误"
	case ERROR_AUTH_NO_TOKEN:
		res = "未提供access_token"
	case ERROR_AUTH_CHECK_TOKEN_FAIL:
		res = "access_token鉴权失败"
	case ERROR_AUTH_CHECK_TOKEN_TIMEOUT:
		res = "access_token已超时"
	case ERROR_AUTH_TOKEN_GENERATE:
		res = "access_token生成失败"
	case ERROR_AUTH:
		res = "access_token错误"
	case ERROR_SUPER_AUTH:
		res = "super_token错误"
	case ERROR_SUPER_AUTH_NO_TOKEN:
		res = "super_token未提供"
	default:
		log.Fatalf("未知错误码[%d]", e)
	}
	return
}
