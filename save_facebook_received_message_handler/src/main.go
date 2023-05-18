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

func handle(ctx context.Context, sqsEvent events.SQSEvent) {
	log.Println("facebook database handler")
	for _, record := range sqsEvent.Records {
		log.Println(record)
	}
	return
}
