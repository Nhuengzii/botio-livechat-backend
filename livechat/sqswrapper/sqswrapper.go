package sqswrapper

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Client struct {
	client *sqs.SQS
}

func NewClient(awsRegion string) *Client {
	sess := session.Must(session.NewSession())
	client := sqs.New(sess, aws.NewConfig().WithRegion(awsRegion))
	return &Client{client}
}

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
