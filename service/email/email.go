package email

import (
	"bytes"
	"delivery-backend/internal/setting"
	"delivery-backend/service/consumer"
	"encoding/json"
	"html/template"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type MsgData struct {
	PageData PageData
	Email    string
}

type PageData struct {
	Account     string
	Password    string
	Name        string
	PhoneNumber string
}

// singleton
type ApprovalSender struct {
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
var sender *ApprovalSender

// singleton
func newApprovalSender() *ApprovalSender {
	f := setting.EmailSetting.TemplatePath
	tmpl, err := template.ParseFiles(f)
	if err != nil {
		log.Panic(err)
	}

	res := ApprovalSender{
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

func (a *ApprovalSender) Connect() error {
	var err error
	a.conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	return err
}

func (a *ApprovalSender) Close() error {
	err := a.conn.Close()
	return err
}

func (a *ApprovalSender) ConsumeMsg() error {
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

			data := MsgData{}
			err := json.Unmarshal(d.Body, &data)
			if err != nil {
				log.Warn(err)
				err = d.Nack(false, false)
				// 这一步是不应该出错的
				log.Panic("UACK error:", err)
			}

			// 发送email
			log.Debugf("Receive email msg[%v]", data)
			err = a.SendMsg(data.Email, &data.PageData)
			if err != nil {
				// 发送失败, NACK
				err = d.Nack(false, false)
				// 这一步是不应该出错的
				log.Panic("UACK error:", err)
				log.Warn(err)
			}

			// 发送成功, ack成功处理
			err = d.Ack(false)
			if err != nil {
				log.Panic("ACK error:", err)
			}
		}
	}()

	return nil
}

func (a *ApprovalSender) SendMsg(to string, data *PageData) error {
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

	/////修改application的状态

	return err
}
