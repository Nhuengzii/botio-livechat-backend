package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

func publishSNSMessage(message string) error {
	sess := session.Must(session.NewSession())
	svc := sns.New(sess, aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")))
	input := &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(snsTopicARN),
	}
	_, err := svc.Publish(input)
	return err
}
