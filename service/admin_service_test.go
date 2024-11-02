package service

import (
	jwt_token "delivery-backend/pkg/jwt"
	"fmt"
	"testing"

	log "github.com/sirupsen/logrus"
)

var (
	account  = "test"
	password = "123456"
)

func TestGetAdminAccessToken(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	tks := jwt_token.GetAccessToken("admin", account, 10)
	fmt.Println(tks)
	account, code := AuthAdminAccessToken(tks)
	fmt.Println(account, code)
}
