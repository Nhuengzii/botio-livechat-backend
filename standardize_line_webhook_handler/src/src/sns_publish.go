package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

var snsTopicArn = os.Getenv("SNS_TOPIC_ARN")

func publishToSNS(message string) error {
	sess := session.Must(session.NewSession())
	svc := sns.New(sess, aws.NewConfig().WithRegion("ap-southeast-1"))
	params := &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(snsTopicArn),
	}
	_, err := svc.Publish(params)
	return err
}
