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
			log.Println("save message handler: " + err.Error())
			logToDiscord("save message handler: " + err.Error())
		}
	}()
	dbc, err := newDBclient(ctx)
	if err != nil {
		return err
	}
	for _, sqsMessage := range sqsEvent.Records {
		snsMessageString := sqsMessage.Body
		snsMsg := &snsMessage{}
		if err := json.Unmarshal([]byte(snsMessageString), snsMsg); err != nil {
			return err
		}
		message := &botioMessage{}
		if err := json.Unmarshal([]byte(snsMsg.Message), message); err != nil {
			return err
		}
		if err := messageHandler(ctx, dbc, message); err != nil {
			return err
		}
	}
	if err := dbc.close(ctx); err != nil {
		return err
	}
	return nil
}
