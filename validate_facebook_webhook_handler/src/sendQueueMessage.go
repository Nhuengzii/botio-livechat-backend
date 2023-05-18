package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func SendQueueMessage(message string, sqsClient *sqs.SQS, queueUrl string) error {
	log.Println("SendQueueMessage Function")
	output, err := sqsClient.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(0),
		QueueUrl:     &queueUrl,
		MessageBody:  &message,
	})

	log.Printf("SQS send message output : %v\n", output)
	return err
}
