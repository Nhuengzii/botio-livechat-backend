package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, sqsEvent events.SQSEvent) (err error) {
	defer func() {
		if err != nil {
			log.Println("message standardizer handler: " + err.Error())
			logToDiscord("message standardizer handler: " + err.Error())
		}
	}()
	for _, sqsMessage := range sqsEvent.Records {
		webhookBodyString := sqsMessage.Body
		wb, err := parseWebhookBody(webhookBodyString)
		if err != nil {
			return err
		}
		botioMessages := wb.toBotioMessages()
		for _, message := range botioMessages {
			messageJSON, err := json.Marshal(message)
			if err != nil {
				return err
			}
			if err := publishSNSMessage(string(messageJSON)); err != nil {
				return err
			}
		}
	}
	return nil
}
