package livechat

// PushNotificationClient is an interface for push notification client, for example, SNS
type PushNotificationClient interface {
	PublishMessage(topicARN string, message string) error
}
