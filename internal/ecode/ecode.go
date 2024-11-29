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

	// wx related error
	ERROR_WX_SESSION_EXPIRE Ecode = 9000
	ERROR_WX_CART_EMPTY     Ecode = 9010
	ERROR_WX_ORDER_CREATE   Ecode = 9011
	ERROR_WX_ORDER_CANCEL   Ecode = 9012
	ERROR_WX_IMAGE_UPLOAD   Ecode = 9013

	// admin related error
	ERROR_ADMIN_NOT_FOUND         Ecode = 10000
	ERROR_ADMIN_INCORRECT_PWD     Ecode = 10001
	ERROR_ADMIN_ACCOUNT_EXIST     Ecode = 10002
	ERROR_ADMIN_LOGOUT            Ecode = 10003
	ERROR_ADMIN_NOT_LOGIN         Ecode = 10004
	ERROR_ADMIN_ROLE              Ecode = 10005
	ERROR_ADMIN_INVALID_OPERATION Ecode = 10006

	// merchant related
	ERROR_MERCHANT_NON_FOUND      Ecode = 11000
	ERROR_MERCHANT_INCORRECT_PWD  Ecode = 11001
	ERROR_MERCHANT_ACCOUNT_EXIST  Ecode = 11002
	ERROR_MERCHANT_LOGOUT         Ecode = 11003
	ERROR_MERCHANT_NOT_LOGIN      Ecode = 11004
	ERROR_MERCHANT_ROLE           Ecode = 11005
	ERROR_MERCHANT_ACCOUNT_BANNED Ecode = 11006
	ERROR_MERCHANT_UNAUTH         Ecode = 11007

	ERROR_MERCHANT_APPLICATION_NOT_FOUND Ecode = 12000

	ERROR_RESTAURANT_EXIST     Ecode = 13000
	ERROR_RESTAURANT_NOT_FOUND Ecode = 13001

	ERROR_CATEGORY_NOT_FOUND Ecode = 14000
	ERROR_DISH_NOT_FOUND     Ecode = 15000
	ERROR_FLAVOR_NOT_FOUND   Ecode = 16000

	ERROR_AUTH_NO_ACCESS_TOKEN       Ecode = 20000
	ERROR_AUTH_ACCESS_TOKEN_EXPIRED  Ecode = 20001
	ERROR_AUTH_NO_REFRESH_TOKEN      Ecode = 20002
	ERROR_AUTH_REFRESH_TOKEN_EXPIRED Ecode = 20003
	ERROR_AUTH_CHECK_ACCESS_TOKEN    Ecode = 20004
	ERROR_AUTH_CHECK_REFRESH_TOKEN   Ecode = 20005

	ERROR_AUTH_TOKEN_GENERATE Ecode = 20010

	ERROR_SUPER_AUTH          Ecode = 30000
	ERROR_SUPER_AUTH_NO_TOKEN Ecode = 30001
)

func StatusText(e Ecode) (res string) {
	switch e {
	case SUCCESS:
		res = "ok"
	case ERROR:
		res = "fail"

	case ERROR_WX_SESSION_EXPIRE:
		res = "微信session过期"
	case ERROR_WX_CART_EMPTY:
		res = "购物车为空，无法下单"
	case ERROR_WX_ORDER_CREATE:
		res = "下单失败"
	case ERROR_WX_ORDER_CANCEL:
		res = "取消订单失败"
	case ERROR_WX_IMAGE_UPLOAD:
		res = "微信图片上传失败"

	case INVALID_PARAMS:
		res = "请求参数错误"
	case ERROR_ADMIN_NOT_FOUND:
		res = "管理员账号不存在"
	case ERROR_ADMIN_ACCOUNT_EXIST:
		res = "管理员账号已注册"
	case ERROR_ADMIN_INCORRECT_PWD:
		res = "管理员密码输入错误"
	case ERROR_ADMIN_LOGOUT:
		res = "管理员非法登出请求/已登出"
	case ERROR_ADMIN_NOT_LOGIN:
		res = "管理员未登录"
	case ERROR_ADMIN_ROLE:
		res = "管理员身份错误"
	case ERROR_ADMIN_INVALID_OPERATION:
		res = "管理员非法操作"

	case ERROR_MERCHANT_NON_FOUND:
		res = "商家账号不存在"
	case ERROR_MERCHANT_ACCOUNT_EXIST:
		res = "商家账号已注册"
	case ERROR_MERCHANT_INCORRECT_PWD:
		res = "商家密码输入错误"
	case ERROR_MERCHANT_LOGOUT:
		res = "商家非法登出请求/已登出"
	case ERROR_MERCHANT_NOT_LOGIN:
		res = "商家未登录"
	case ERROR_MERCHANT_ROLE:
		res = "商家身份错误"
	case ERROR_MERCHANT_ACCOUNT_BANNED:
		res = "商家账号被禁用"
	case ERROR_MERCHANT_UNAUTH:
		res = "商家对请求数据没有权限"

	case ERROR_MERCHANT_APPLICATION_NOT_FOUND:
		res = "商家申请表记录不存在"

	case ERROR_RESTAURANT_EXIST:
		res = "该商铺名已经存在"
	case ERROR_RESTAURANT_NOT_FOUND:
		res = "商铺不存在"

	case ERROR_CATEGORY_NOT_FOUND:
		res = "菜品分类不存在"

	case ERROR_DISH_NOT_FOUND:
		res = "菜品不存在"

	case ERROR_FLAVOR_NOT_FOUND:
		res = "口味不存在"

		// refresh_token and access_token
	case ERROR_AUTH_NO_REFRESH_TOKEN:
		res = "未提供refresh_token"
	case ERROR_AUTH_CHECK_REFRESH_TOKEN:
		res = "refresh_token鉴权失败"
	case ERROR_AUTH_REFRESH_TOKEN_EXPIRED:
		res = "refresh_token过期"
	case ERROR_AUTH_NO_ACCESS_TOKEN:
		res = "未提供access_token"
	case ERROR_AUTH_CHECK_ACCESS_TOKEN:
		res = "access_token鉴权失败"
	case ERROR_AUTH_ACCESS_TOKEN_EXPIRED:
		res = "access_token过期"

		// token generate
	case ERROR_AUTH_TOKEN_GENERATE:
		res = "access_token生成失败"
	case ERROR_SUPER_AUTH:
		res = "super_token错误"
	case ERROR_SUPER_AUTH_NO_TOKEN:
		res = "super_token未提供"
	default:
		log.Panicf("未知错误码[%d]", e)
	}
	return
}
