package producer

type Producer interface {
	Connect(url string) error
	PublishMsg(queueName string, msg any) error
	Close() error
}
