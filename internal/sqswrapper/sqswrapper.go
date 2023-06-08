package sqswrapper

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Client struct {
	client *sqs.SQS
}

func NewClient() *Client {
	sess := session.Must(session.NewSession())
	client := sqs.New(sess, aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")))
	return &Client{client}
}

func (c *Client) SendMessage(queueURL string, message string) error {
	input := &sqs.SendMessageInput{
		MessageBody: aws.String(message),
		QueueUrl:    aws.String(queueURL),
	}
	_, err := c.client.SendMessage(input)
	if err != nil {
		return fmt.Errorf("sqs.SendMessage: %w", err)
	}
	return nil
}
