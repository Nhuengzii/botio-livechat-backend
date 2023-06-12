package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/snswrapper"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	errNoMessageEntry       = errors.New("Error! no message entry")
	errUnknownWebhookType   = errors.New("Error! unknown webhook type found!")
	errUnknownWebhookObject = errors.New("Error! unknown webhook Object found!")
)

func main() {
	c := config{
		DiscordWebhookURL:       os.Getenv("DISCORD_WEBHOOK_URL"),
		SnsTopicARN:             os.Getenv("SNS_TOPIC_ARN"),
		SnsClient:               snswrapper.NewClient(os.Getenv("AWS_REGION")),
		FacebookPageAccessToken: os.Getenv("ACCESS_TOKEN"),
	}
	lambda.Start(c.handler)
}

func (c *config) handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	discord.Log(c.DiscordWebhookURL, "facebook standardize webhook handler")
	start := time.Now()
	var recieveWebhook ReceiveWebhook
	for _, record := range sqsEvent.Records {
		err := json.Unmarshal([]byte(record.Body), &recieveWebhook)
		if err != nil {
			discord.Log(c.DiscordWebhookURL, fmt.Sprintf("error unmarshal recieve webhook "))
			return errUnknownWebhookObject
		}
		c.handleRecieveWebhook(&recieveWebhook)
	}
	discord.Log(c.DiscordWebhookURL, fmt.Sprintf("Elapsed: %v", time.Since(start)))
	return nil
}
