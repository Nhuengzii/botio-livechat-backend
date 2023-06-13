package livechat

type PushNotificationClient interface {
	PublishMessage(topicARN string, message string) error
}
