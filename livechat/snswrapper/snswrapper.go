package snswrapper

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type Client struct {
	client *sns.SNS
}

func NewClient(awsRegion string) *Client {
	sess := session.Must(session.NewSession())
	client := sns.New(sess, aws.NewConfig().WithRegion(awsRegion))
	return &Client{client}
}

func (c *Client) PublishMessage(topicARN string, v any) error {
	message, err := json.Marshal(v)
	if err != nil {
		return err
	}
	snsPublishMessage := SNSPublishMessage{
		Default: string(message),
	}
	snsByte, err := json.Marshal(snsPublishMessage)
	if err != nil {
		return err
	}

	input := &sns.PublishInput{
		Message:  aws.String(string(snsByte)),
		TopicArn: aws.String(topicARN),
	}
	_, err = c.client.Publish(input)
	if err != nil {
		return fmt.Errorf("sns.PublishMessage: %w", err)
	}
	return nil
}

type SNSPublishMessage struct {
	Default string `json:"default"`
}

type SNSMessage struct {
	Message string `json:"message"`
}
