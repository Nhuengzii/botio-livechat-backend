package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func sendSQSMessage(message string) error {
	sess := session.Must(session.NewSession())
	svc := sqs.New(sess, aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")))
	input := &sqs.SendMessageInput{
		MessageBody: aws.String(message),
		QueueUrl:    aws.String(sqsQueueURL),
	}
	_, err := svc.SendMessage(input)
	if err != nil {
		return &sendSQSMessageError{
			message: "couldn't send message to sqs",
			err:     err,
		}
	}
	return nil
}

type sendSQSMessageError struct {
	message string
	err     error
}

func (e *sendSQSMessageError) Error() string {
	return e.message + ": " + e.err.Error()
}
