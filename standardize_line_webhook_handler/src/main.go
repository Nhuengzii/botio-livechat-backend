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

func Handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	log.Println("message standardizer handler: received sqs event")
	logToDiscord("message standardizer handler: received sqs event")
	for _, sqsMessage := range sqsEvent.Records {
		webhookBodyString := sqsMessage.Body
		wb, err := parseWebhookBody(webhookBodyString)
		if err != nil {
			log.Println("message standardizer handler: " + err.Error())
			logToDiscord("message standardizer handler: " + err.Error())
			return err
		}
		botioMessages := wb.toBotioMessages()
		for _, message := range botioMessages {
			messageJSON, err := json.Marshal(message)
			if err != nil {
				log.Println("message standardizer handler: couldn't marshal botio message: " + err.Error())
				logToDiscord("message standardizer handler: couldn't marshal botio message: " + err.Error())
				return err
			}
			if err := publishSNSMessage(string(messageJSON)); err != nil {
				log.Println("message standardizer handler: " + err.Error())
				logToDiscord("message standardizer handler: " + err.Error())
				return err
			}
		}
	}
	log.Println("message standardizer handler: published all botio messages to sns")
	logToDiscord("message standardizer handler: published all botio messages to sns")
	return nil
}
