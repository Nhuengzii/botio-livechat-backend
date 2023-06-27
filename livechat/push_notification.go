package livechat

// PushNotificationClient is an interface for push notification client, for example, SNS
type PushNotificationClient interface {
	// PublishMessage recieve a message string and publish the message into specific topic or channel.
	// Return an error if it occurs.
	//
	// If implements by SNS services topicAddress should be topicARN.
	PublishMessage(topicAddress string, message string) error
}
