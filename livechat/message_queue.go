package livechat

// MessageQueueClient is an interface for message queue client, for example, SQS
type MessageQueueClient interface {
	SendMessage(queueURL string, message string) error
}
