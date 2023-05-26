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
	if err != nil {
		return &publishSNSMessageError{
			message: "couldn't publish sns message",
			err:     err,
		}
	}
	return nil
}

type publishSNSMessageError struct {
	message string
	err     error
}

func (e *publishSNSMessageError) Error() string {
	return e.message + ": " + e.err.Error()
}
