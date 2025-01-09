package api_test

import (
	"bytes"
	"delivery-backend/common/util"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"testing"
)

func createAndLogin(account string, pwd string, admin_name string) error {
	url := "http://127.0.0.1:8000/api/v1/admin/create?password=%s&admin_name=%s&account=%s"
	url = fmt.Sprintf(url, pwd, admin_name, account)
	resp, err := http.Post(url, "multipart/form-data", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	log.Println("create status: ", resp.Status)
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println("create resp: ", string(res))

	// 登陆
	url = "http://localhost:8000/api/v1/admin/login"
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	_ = writer.WriteField("account", account)
	_ = writer.WriteField("password", pwd)
	err = writer.Close()
	if err != nil {
		return err
	}
	resp, err = http.Post(
		url,
		writer.FormDataContentType(),
		&requestBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	log.Println("login status: ", resp.Status)
	res, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println("login resp: ", string(res))

	log.Println("create and login once")
	return nil
}

// 基于api层面的测试
func TestAdminLogin(t *testing.T) {
	// 首先测试admin create

	for i := 0; i < 10; i++ {
		account := util.RandString(15)
		pwd := util.RandString(16)
		admin_name := util.RandString(10)
		err := createAndLogin(account, pwd, admin_name)
		if err != nil {
			t.Error(err)
			return
		}
	}
}
