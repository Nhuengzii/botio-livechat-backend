package livechat

type SNSClient interface {
	PublishMessage(topicARN string, v any) error
}
