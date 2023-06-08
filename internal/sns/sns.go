package sns

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type Client struct {
	client *sns.SNS
}

func NewClient() *Client {
	sess := session.Must(session.NewSession())
	client := sns.New(sess, aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")))
	return &Client{client}
}

func (c *Client) PublishMessage(topicARN string, message string) error {
	input := &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(topicARN),
	}
	_, err := c.client.Publish(input)
	if err != nil {
		return fmt.Errorf("sns.PublishMessage: %w", err)
	}
	return nil
}

type SNSMessage struct {
	Message string `json:"message"`
}
