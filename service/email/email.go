package email

import (
	"bytes"
	"delivery-backend/internal/setting"
	"delivery-backend/models"
	"delivery-backend/pkg/utils"
	"delivery-backend/service/consumer"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type MsgData struct {
	ApplicationID uint
	PhoneNumber   string
	Email         string
	Name          string
}

type PageData struct {
	Account     string
	Password    string
	Name        string
	PhoneNumber string
}

// singleton
type EmailSender struct {
	tmpl *template.Template
	from string
	pwd  string
	cc   string
	host string
	port int

	conn      *amqp.Connection
	queueName string
}

func Setup() {
	sender = newApprovalSender()
	err := consumer.Register("email sender", sender)
	if err != nil {
		log.Panic(err)
	}
}

// singleton
var sender *EmailSender

// singleton
func newApprovalSender() *EmailSender {
	f := setting.EmailSetting.TemplatePath
	tmpl, err := template.ParseFiles(f)
	if err != nil {
		log.Panic(err)
	}

	res := EmailSender{
		tmpl:      tmpl,
		from:      setting.EmailSetting.SenderEmail,
		pwd:       setting.EmailSetting.SenderPassword,
		cc:        setting.EmailSetting.CCEmail,
		host:      setting.EmailSetting.SMTPHost,
		port:      setting.EmailSetting.SMTPPort,
		queueName: "merchant_approval_email",
	}

	return &res
}

func (a *EmailSender) Connect() error {
	var err error
	a.conn, err = amqp.Dial(setting.RabitmqSetting.DialURL)
	return err
}

func (a *EmailSender) Close() error {
	err := a.conn.Close()
	return err
}

func (a *EmailSender) ConsumeMsg() error {
	ch, err := a.conn.Channel()
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		a.queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	// consume消息
	msgs, err := ch.Consume(
		q.Name,         // queue
		"email sender", // consumer
		false,          // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		return err
	}

	// 执行email的事务处理
	go func() {
		log.Info("email sender waiting for message")

		for d := range msgs {
			log.Trace("Receive approve msg")
			err := a.approveHandler(&d)
			if err != nil {
				// send NACK
				log.Warn(err)
				err = d.Nack(false, false)
				if err != nil {
					log.Panic("UACK error:", err)
				}
			} else {
				// 成功处理, send ACK
				err = d.Ack(false)
				if err != nil {
					log.Panic("ACK error:", err)
				}
			}
		}
	}()

	return nil
}

// 1. 接收消息
// 2. 创建商家账户
// 3. 发送邮件
// 4. 修改申请表发送邮件的状态
// * ACK/NACK
func (a *EmailSender) approveHandler(d *amqp.Delivery) (err error) {
	///////解析data
	data := MsgData{}
	err = json.Unmarshal(d.Body, &data)
	if err != nil {
		return err
	}

	/////增添商家账户
	// TODO: 暂定注册规则为随机字符串,后续按照需要更改
	account := "M" + utils.RandString(10)
	password := "P" + utils.RandString(11)
	en_password := utils.Encrypt(password, setting.AppSetting.Salt)
	log.Trace("Merchant Account/PWD generated")

	m := models.Merchant{
		Account:               account,
		Password:              en_password,
		PhoneNumber:           data.PhoneNumber,
		MerchantName:          data.Name,
		MerchantApplicationID: data.ApplicationID,
	}
	created, err := models.CreateMerchant(&m)
	if err != nil {
		return err
	}
	if !created {
		// 说明之前已经handle过这个request了, 就不必继续处理, 重复发送了create请求
		// 直接标记处理完毕即可
		return fmt.Errorf("merchant with application_id[%v] exists", m.MerchantApplicationID)
	}
	log.Tracef("created merchant: ACCOUNT[%s],PWD[%s]", account, password)

	/////最后才是发送邮件
	pageData := PageData{
		Account:     account,
		Password:    password,
		PhoneNumber: data.PhoneNumber,
		Name:        data.Name,
	}

	err = a.sendEmail(data.Email, &pageData)
	if err != nil {
		m_err := models.UpdateEmailStatus(data.ApplicationID, models.EmailSentError)
		if m_err != nil {
			log.Warn(m_err)
		}
		return err
	}

	/////修改application中的邮件的状态
	err = models.UpdateEmailStatus(data.ApplicationID, models.EmailSent)
	if err != nil {
		log.Warn(err)
		return
	}
	log.Trace("Application email status set to EmailSent")

	return nil
}

func (a *EmailSender) sendEmail(to string, data *PageData) error {
	var mailbody bytes.Buffer
	err := a.tmpl.Execute(&mailbody, data)
	if err != nil {
		return err
	}

	if !setting.EmailSetting.EmailOn {
		// 没有开启email
		log.Info(
			"email mode off\n",
			"simulating sent email...")
		time.Sleep(time.Duration(1) * time.Second)
		f, err := os.Create(data.Account + "-email.html")
		if err != nil {
			log.Warn(err)
			log.Info(mailbody.String())
		}
		// 写入文件
		f.Write(mailbody.Bytes())
		log.Infof("email sent from[%v] to [%v]",
			a.from, to)
		return nil
	}

	/////准备邮件的发送
	m := gomail.NewMessage()
	m.SetHeader("From", a.from)
	m.SetHeader("To" /*"2755345380@qq.com",*/, to)
	m.SetHeader("Subject", "您的商家申请已通过")
	if a.cc != "" {
		m.SetAddressHeader("Cc", a.cc, "admin")
	}

	m.SetBody("text/html", mailbody.String())
	d := gomail.NewDialer(
		a.host,
		a.port,
		a.from,
		a.pwd,
	)
	err = d.DialAndSend(m)
	log.Tracef("send mail with error[%v]", err)
	if err != nil {
		return err
	}

	return err
}
