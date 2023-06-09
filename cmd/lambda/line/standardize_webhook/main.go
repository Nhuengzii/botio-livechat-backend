package main

import (
	"context"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/snswrapper"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"os"
)

func (c *config) handler(ctx context.Context, sqsEvent events.SQSEvent) (err error) {
	defer func() {
		if err != nil {
			logMessage := "lambda/line/standardize_webhook/main.config.handler: " + err.Error()
			log.Println(logMessage)
			discord.Log(c.discordWebhookURL, logMessage)
		}
	}()
	for _, sqsMessage := range sqsEvent.Records {
		hookBody, err := parseWebhookBody(sqsMessage.Body)
		if err != nil {
			return err
		}
		err = handleEvents(c, hookBody)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	c := &config{
		discordWebhookURL: os.Getenv("DISCORD_WEBHOOK_URL"),
		snsTopicARN:       os.Getenv("SNS_TOPIC_ARN"),
		snsClient:         snswrapper.NewClient(os.Getenv("AWS_REGION")),
	}
	lambda.Start(c.handler)
}
