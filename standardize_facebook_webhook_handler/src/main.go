package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, sqsEvent events.SQSEvent) error {
	log.Println("Facebook Message Standardizer handler")
	for _, message := range sqsEvent.Records {
		log.Printf("The message body is : %s\n", message.Body)
	}
	return nil
}
