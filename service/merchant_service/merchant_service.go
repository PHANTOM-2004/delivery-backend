package merchant_service

import (
	"context"
	"delivery-backend/internal/app"
	"delivery-backend/internal/ecode"
	"delivery-backend/internal/setting"
	"delivery-backend/models"
	"delivery-backend/pkg/utils"
	"delivery-backend/service/email"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type Password struct {
	Password string
}

type Login struct {
	Account  string
	Password string
}

type SignUp struct {
	Account      string
	Password     string
	MerchantName string
	PhoneNumber  string
}

func init() {
	// NOTE: 在该init中，初始化该模块的数据验证
	{
		// register login validation
		// 修改密码validation
		password_rules := map[string]string{
			"Password": "min=8,max=30",
		}

		app.RegisterMapValidation(Password{}, password_rules)

		// 登录validation
		login_rules := password_rules
		login_rules["Account"] = "min=6,max=30"

		app.RegisterMapValidation(Login{}, login_rules)

		// 注册账号validation
		signup_rules := login_rules
		signup_rules["MerchantName"] = "min=2,max=20"
		// example: +8613912345678
		signup_rules["PhoneNumber"] = "required,e164"
		app.RegisterMapValidation(SignUp{}, signup_rules)
	}
}

// 如果认证成功，会返回id
func MerchantLoginValidate(account string, password string, c *gin.Context) (uint, bool) {
	data := Login{
		Account:  account,
		Password: password,
	}

	err := app.ValidateStruct(&data)
	if err != nil {
		// 通常来说前端不应当传递非法参数，对于非法参数的传递
		// 通常是其他人所进行的
		log.Warn("Login: invalid params")
		log.Debug(err, data)
		app.ResponseInvalidParams(c)
		return 0, false
	}

	m, err := models.GetMerchant(data.Account)
	if err != nil {
		// 其他未知错误
		log.Warn(err)
		app.ResponseInternalError(c, err)
		return 0, false
	}
	if m == nil {
		// 商家不存在
		log.Debug(err, data)
		app.Response(c, http.StatusOK, ecode.ERROR_MERCHANT_NON_FOUND, nil)
		return 0, false
	}
	if m.Status == models.MerchantAccountDisabled {
		// 商家账号被禁用
		app.Response(c, http.StatusOK, ecode.ERROR_MERCHANT_ACCOUNT_BANNED, nil)
		return 0, false
	}

	en_pwd := utils.Encrypt(data.Password, setting.AppSetting.Salt)
	if en_pwd != m.Password {
		// 用户输错密码
		log.Debug("incorrect password")
		app.Response(c, http.StatusOK, ecode.ERROR_MERCHANT_INCORRECT_PWD, nil)
		return 0, false
	}
	return m.ID, true
}

func SignUpRequestValidate(c *gin.Context) bool {
	method := c.Request.Method
	if method != "POST" {
		log.Fatal("invalid usage")
		return false
	}

	data := SignUp{
		Account:      c.PostForm("account"),
		Password:     c.PostForm("password"),
		MerchantName: c.PostForm("merchant_name"),
		PhoneNumber:  c.PostForm("phone_number"),
	}

	err := app.ValidateStruct(&data)
	if err != nil {
		app.ResponseInvalidParams(c)
		log.Warn(err)
		return false
	}
	return true
}

func PasswordRequestValidate(password string, c *gin.Context) bool {
	data := Password{password}

	err := app.ValidateStruct(&data)
	if err != nil {
		app.ResponseInvalidParams(c)
		log.Warn(err)
		return false
	}
	return true
}

// queue_name: merchant_approval_email
type EmailProducer struct {
	queueName string
	// exchangeName string
	conn *amqp.Connection
}

func NewEmailProducer() *EmailProducer {
	res := EmailProducer{
		queueName: "merchant_approval_email",
		// exchangeName: "merchant_email",
	}
	return &res
}

func (e *EmailProducer) Connect() error {
	var err error
	e.conn, err = amqp.Dial(setting.RabitmqSetting.DialURL)
	return err
}

func (e *EmailProducer) Close() error {
	err := e.conn.Close()
	return err
}

func (e *EmailProducer) PublishMsg(data *email.MsgData) error {
	ch, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		e.queueName, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)
	if err == nil {
		log.Tracef("publish email msg [%v]", string(body))
	}
	return err
}

// 管理员端为Merchant创建
func CreateMerchantFromApplication(application_id uint) error {
	a, err := models.
		GetMerchantApplication(application_id)
	if err != nil {
		return err
	} else if a.ID == 0 {
		return fmt.Errorf("merchant application [%v] not found", application_id)
	}

	data := email.MsgData{
		ApplicationID: application_id,
		PhoneNumber:   a.PhoneNumber,
		Email:         a.Email,
		Name:          a.Name,
	}

	// 接下来发送email, 作为生产者推送给email 服务
	p := NewEmailProducer()
	err = p.Connect()
	if err != nil {
		return err
	}
	defer func() {
		err := p.Close()
		if err != nil {
			log.Warn(err)
		}
	}()
	err = p.PublishMsg(&data)
	if err == nil {
		log.Tracef("email to merchant[%v] sent", data.Name)
	}
	return err
}
