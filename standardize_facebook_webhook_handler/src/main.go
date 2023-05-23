package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	errNoMessageEntry     = errors.New("Error! no message entry")
	errUnknownWebhookType = errors.New("Error! unknown webhook type found!")
)

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, sqsEvent events.SQSEvent) {
	discordLog(fmt.Sprintln("Facebook Message Standardizer handler"))
	start := time.Now()
	var recieveMessage RecieveMessage
	for _, record := range sqsEvent.Records {
		err := json.Unmarshal([]byte(record.Body), &recieveMessage)
		if err != nil || recieveMessage.Object != "page" {
			discordLog(fmt.Sprintf("Error unknown webhook object: %v\n", err))
			return
		}
		log.Printf("%+v", recieveMessage)
		for _, message := range recieveMessage.Entry {
			err = handleWebhookEntry(message)
			if err != nil {
				discordLog(fmt.Sprintf("Error handling webhook entry : %v", err))
			}
		}
	}
	discordLog(fmt.Sprintf("Elapsed: %v", time.Since(start)))
	return
}

func handleWebhookEntry(message Notification) error {
	if len(message.MessageDatas) <= 0 {
		return errNoMessageEntry
	}

	for _, messageData := range message.MessageDatas {
		if messageData.Message.MessageID != "" {
			// standardize messaging hooks
			var standardMessage StandardMessage
			err := StandardizeMessage(messageData, message.PageID, &standardMessage)
			if err != nil {
				return err
			}
			err = sendSnsMessage(&standardMessage)
			if err != nil {
				return err
			}
		} else if messageData.Delivery.Watermark != 0 {
			// standardize delivery hooks
			discordLog(fmt.Sprintf("Delivery Webhook"))
		} else {
			return errUnknownWebhookType
		}
	}
	return nil
}
