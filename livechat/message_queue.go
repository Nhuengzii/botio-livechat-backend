package livechat

// MessageQueueClient is an interface for message queue client, for example, SQS
type MessageQueueClient interface {
	// SendMessage recieve a message string and send the message into specific message queue.
	// Return an error if it occurs.
	SendMessage(queueURL string, message string) error
}
