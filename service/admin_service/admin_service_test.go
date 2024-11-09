package admin_service

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

	tks := jwt_token.GetAccessToken("admin", 10, account, 10)
	fmt.Println(tks)
	id, account, code := AuthAccessToken(tks)
	fmt.Println(id, account, code)
}
