package main

import (
	"fmt"
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
	if _, err := svc.SendMessage(input); err != nil {
		return fmt.Errorf("sendSQSMessage: %w", err)
	}
	return nil
}
