package livechat

// MessageQueueClient is an interface for message queue client, for example, SQS
type MessageQueueClient interface {
	// SendMessage recieve a message string and send the message into specific message queue.
	// Return an error if it occurs.
	//
	// If implements by SQS services queueAddress should be queueURL.
	SendMessage(queueAddress string, message string) error
}
