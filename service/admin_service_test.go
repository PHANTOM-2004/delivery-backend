package service

import (
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

	tks := GetAdminAccessToken(account)
	fmt.Println(tks)
	account, code := AuthAdminAccessToken(tks)
	fmt.Println(account, code)
}
