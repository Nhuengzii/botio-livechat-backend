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
	log.Println("from standardizer: received sqs event")
	if logToDiscordEnabled {
		logToDiscord("from standardizer: received sqs event")
	}
	for _, sqsMessage := range sqsEvent.Records {
		wBody := sqsMessage.Body
		wb, err := parseWebhookBody(wBody)
		if err != nil {
			log.Println("from standardizer: couldn't parse webhook body")
			if logToDiscordEnabled {
				logToDiscord("from standardizer: couldn't parse webhook body")
			}
			return err
		}
		botioMessages := wb.toBotioMessages()
		for _, message := range botioMessages {
			messageJSON, err := json.Marshal(message)
			if err != nil {
				log.Println("from standardizer: couldn't marshal botio message")
				if logToDiscordEnabled {
					logToDiscord("from standardizer: couldn't marshal botio message")
				}
				return err
			}
			if err := publishSNSMessage(string(messageJSON)); err != nil {
				log.Println("from standardizer: couldn't publish botio message")
				if logToDiscordEnabled {
					logToDiscord("from standardizer: couldn't publish botio message")
				}
				return err
			}
		}
	}
	log.Println("from standardizer: published all botio messages to sns")
	if logToDiscordEnabled {
		logToDiscord("from standardizer: published all botio messages to sns")
	}
	return nil
}
