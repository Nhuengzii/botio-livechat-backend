package main

import (
	"context"
	"os"

	"github.com/Nhuengzii/botio-livechat-backend/internal/sqswrapper"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Lambda struct {
	config
}

func (l Lambda) handler(ctx context.Context, event events.SQSEvent) {
}

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
