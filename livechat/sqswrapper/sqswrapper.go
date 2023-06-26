// Package sqswrapper implements SQS's Client for manipulating SQS's service database
package sqswrapper

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// A Client contains SQS's client and a Target struct.
type Client struct {
	client *sqs.SQS // SQS's client used to do various SQS's operation
}

// NewClient create and return a SQS's client.
func NewClient(awsRegion string) *Client {
	sess := session.Must(session.NewSession())
	client := sqs.New(sess, aws.NewConfig().WithRegion(awsRegion))
	return &Client{client}
}

// SendMessage recieve a message string and send the message into specific SQS queue.
// Return an error if it occurs.
func (c *Client) SendMessage(queueURL string, message string) error {
	input := &sqs.SendMessageInput{
		MessageBody: aws.String(message),
		QueueUrl:    aws.String(queueURL),
	}
	_, err := c.client.SendMessage(input)
	if err != nil {
		return fmt.Errorf("sqswrapper.SendMessage: %w", err)
	}
	return nil
}
