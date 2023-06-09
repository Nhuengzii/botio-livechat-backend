package livechat

type MessageQueueClient interface {
	SendMessage(queueURL string, message string) error
}
