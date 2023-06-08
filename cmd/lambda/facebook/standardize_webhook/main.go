package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/internal/discord"
	"github.com/Nhuengzii/botio-livechat-backend/internal/fbutil/webhook"
	"github.com/Nhuengzii/botio-livechat-backend/internal/snswrapper"
	"github.com/Nhuengzii/botio-livechat-backend/pkg/stdmessage"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Lambda struct {
	config
}

var (
	errNoMessageEntry       = errors.New("Error! no message entry")
	errUnknownWebhookType   = errors.New("Error! unknown webhook type found!")
	errUnknownWebhookObject = errors.New("Error! unknown webhook Object found!")
)

func main() {
	l := Lambda{
		config: config{
			DiscordWebhookURL:       os.Getenv("DISCORD_WEBHOOK_URL"),
			SnsQueueURL:             os.Getenv("SNS_QUEUE_URL"),
			SnsClient:               *snswrapper.NewClient(),
			FacebookPageAccessToken: os.Getenv("ACCESS_TOKEN"),
		},
	}
	lambda.Start(l.handler)
}

func (l *Lambda) handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	discord.Log(l.DiscordWebhookURL, "facebook standardize webhook handler")
	start := time.Now()
	var recieveMessage webhook.RecieveMessage
	for _, record := range sqsEvent.Records {
		err := json.Unmarshal([]byte(record.Body), &recieveMessage)
		if err != nil || recieveMessage.Object != "page" {
			discord.Log(l.DiscordWebhookURL, fmt.Sprintf("Error unknown webhook object: %v\n", err))
			return errUnknownWebhookObject
		}
		for _, message := range recieveMessage.Entry {
			err = l.handleWebhookEntry(message)
			if err != nil {
				discord.Log(l.DiscordWebhookURL, fmt.Sprintf("Error handling webhook entry : %v", err))
				return err
			}
		}
	}
	discord.Log(l.DiscordWebhookURL, fmt.Sprintf("Elapsed: %v", time.Since(start)))
	return nil
}

func (l *Lambda) handleWebhookEntry(message webhook.Notification) error {
	if len(message.MessageDatas) <= 0 {
		return errNoMessageEntry
	}

	for _, messageData := range message.MessageDatas {
		if messageData.Message.MessageID != "" {
			// standardize messaging hooks
			var standardMessage stdmessage.StdMessage
			err := messageData.StandardizeMessage(l.FacebookPageAccessToken, message.PageID, &standardMessage)
			if err != nil {
				return err
			}

			err = l.SnsClient.PublishMessage(l.SnsQueueURL, standardMessage)
			if err != nil {
				return err
			}
		} else {
			return errUnknownWebhookType
		}
	}
	return nil
}
