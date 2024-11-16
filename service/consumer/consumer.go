package consumer

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Consumer interface {
	Connect() error
	// 持续监听队列, consume消息
	ConsumeMsg() error
	Close() error
}

var consumers = make(map[string]Consumer)

func Register(name string, c Consumer) error {
	if c == nil {
		return fmt.Errorf("consumer[%v] must be initialized", name)
	}
	if cc, ok := consumers[name]; ok {
		return fmt.Errorf("consumer[%v][%v] already registered", name, cc)
	}
	consumers[name] = c
	log.Infof("registered consumer[%v]", name)
	return nil
}

func Setup() {
	for name, c := range consumers {
		log.Infof("Launching consumer[%v]", name)
		err := c.Connect()
		if err != nil {
			log.Fatal(err)
		}

		err = c.ConsumeMsg()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func ShutDown() {
	for name, c := range consumers {
		log.Infof("Shutdown comsumer[%v]", name)
		err := c.Close()
		if err != nil {
			log.Warn(err)
		}
	}
}
