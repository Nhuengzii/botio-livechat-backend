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

func Handler(ctx context.Context, sqsEvent events.SQSEvent) {
	discordLog("Received event")
	// todo: implement deduplication logic
	for _, sqsMessage := range sqsEvent.Records {
		webhookBody := sqsMessage.Body
		discordLog(webhookBody)
		wb, err := parseRequest([]byte(webhookBody))
		if err != nil {
			discordLog("Error parsing webhook body")
			log.Fatal(err)
		}
		discordLog("Parsed webhook body")
		standardMessages := wb.toStandardMessages()
		for _, standardMessage := range standardMessages {
			stdMsg, _ := json.Marshal(standardMessage)
			discordLog(string(stdMsg))
			err := publishToSNS(string(stdMsg))
			if err != nil {
				discordLog("Error publishing to SNS")
				log.Fatal(err)
			}
			discordLog("Published to SNS")
		}
	}
}

func parseRequest(webhookBody []byte) (WebhookBody, error) {
	wb := WebhookBody{}
	if err := json.Unmarshal(webhookBody, &wb); err != nil {
		return wb, err
	}
	return wb, nil
}
