package livechat

type SQSClient interface {
	SendMessage(queueURL string, message string) error
}
