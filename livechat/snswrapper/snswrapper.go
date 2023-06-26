package snswrapper

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

// A Client contains mongodb client and a Target struct.
type Client struct {
	client *sns.SNS // SNS's client used to do various SNS's operation
}

// NewClient create and return a SNS's client.
func NewClient(awsRegion string) *Client {
	sess := session.Must(session.NewSession())
	client := sns.New(sess, aws.NewConfig().WithRegion(awsRegion))
	return &Client{client}
}

// PublishMessage recieve a message string and publish the message into specific SNS topic.
// Return an error if it occurs.
func (c *Client) PublishMessage(topicARN string, message string) error {
	input := &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(topicARN),
	}
	_, err := c.client.Publish(input)
	if err != nil {
		return fmt.Errorf("snswrapper.PublishMessage: %w", err)
	}
	return nil
}

// A SNSMessage contains an actual message recieved via SNS topic.
//
// This struct should be use to recieve and then unmarshal the actual message.
type SNSMessage struct {
	Message string `json:"message"`
}
