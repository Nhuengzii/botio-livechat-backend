package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/internal/discord"
	"github.com/Nhuengzii/botio-livechat-backend/internal/fbutil/standardize"
	"github.com/Nhuengzii/botio-livechat-backend/internal/sqswrapper"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Lambda struct {
	config
}

var (
	errNoMessageEntry     = errors.New("Error! no message entry")
	errUnknownWebhookType = errors.New("Error! unknown webhook type found!")
)

func main() {
	l := Lambda{
		config: config{
			DiscordWebhookURL: os.Getenv("DISCORD_WEBHOOK_URL"),
			SqsQueueURL:       os.Getenv("SQS_QUEUE_URL"),
			SqsClient:         *sqswrapper.NewClient(),
		},
	}
	lambda.Start(l.handler)
}

func (l Lambda) handler(ctx context.Context, sqsEvent events.SQSEvent) {
	discord.Log(l.DiscordWebhookURL, "facebook standardize webhook handler")
	start := time.Now()
	var recieveMessage standardize.RecieveMessage
	for _, record := range sqsEvent.Records {
		err := json.Unmarshal([]byte(record.Body), &recieveMessage)
		if err != nil || recieveMessage.Object != "page" {
			discord.Log(l.DiscordWebhookURL, fmt.Sprintf("Error unknown webhook object: %v\n", err))
			return
		}
		for _, message := range recieveMessage.Entry {
			err = handleWebhookEntry(message)
			if err != nil {
				discord.Log(l.DiscordWebhookURL, fmt.Sprintf("Error handling webhook entry : %v", err))
			}
		}
	}
	discord.Log(l.DiscordWebhookURL, fmt.Sprintf("Elapsed: %v", time.Since(start)))
}

func handleWebhookEntry(message standardize.Notification) error {
	if len(message.MessageDatas) <= 0 {
		return errNoMessageEntry
	}

	for _, messageData := range message.MessageDatas {
		if messageData.Message.MessageID != "" {
			// standardize messaging hooks
			// var standardMessage stdmessage.StdMessage
			// err := StandardizeMessage(messageData, message.PageID, &standardMessage)
			// if err != nil {
			// 	return err
			// }
		} else {
			return errUnknownWebhookType
		}
	}
	return nil
}
