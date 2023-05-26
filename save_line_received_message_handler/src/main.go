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
	log.Println("database write handler: received sqs event")
	logToDiscord("database write handler: received sqs event")
	dbc, err := newDBclient(ctx)
	if err != nil {
		log.Println("database write handler: " + err.Error())
		logToDiscord("database write handler: " + err.Error())
		return err
	}
	for _, sqsMessage := range sqsEvent.Records {
		snsMessageString := sqsMessage.Body
		snsMsg := &snsMessage{}
		if err := json.Unmarshal([]byte(snsMessageString), snsMsg); err != nil {
			log.Println("database write handler: couldn't unmarshal sns message: " + err.Error())
			logToDiscord("database write handler: couldn't unmarshal sns message: " + err.Error())
			return err
		}
		message := &botioMessage{}
		if err := json.Unmarshal([]byte(snsMsg.Message), message); err != nil {
			log.Println("database write handler: couldn't unmarshal botio message: " + err.Error())
			logToDiscord("database write handler: couldn't unmarshal botio message: " + err.Error())
			return err
		}
		if err := messageHandler(ctx, dbc, message); err != nil {
			log.Println("database write handler: " + err.Error())
			logToDiscord("database write handler: " + err.Error())
			return err
		}
	}
	if err := dbc.close(ctx); err != nil {
		log.Println("database write handler: " + err.Error())
		logToDiscord("database write handler: " + err.Error())
		return err
	}
	log.Println("database write handler: wrote all botio messages/conversations to database")
	logToDiscord("database write handler: wrote all botio messages/conversations to database")
	return nil
}
