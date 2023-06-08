package main

import (
	"context"
	"os"

	"github.com/Nhuengzii/botio-livechat-backend/internal/sqswrapper"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Lambda struct {
	config config
}

func (l Lambda) handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{}, nil
}

func main() {
	l := Lambda{
		config: config{
			DiscordWebhookURL: os.Getenv("DISCORD_WEBHOOK_URL"),
			SqsQueueURL:       os.Getenv("SQS_QUEUE_URL"),
			FacebookAppSecret: os.Getenv("FACEBOOK_APP_SECRET"), // TODO to be removed and get from some db instead
			SqsClient:         *sqswrapper.NewClient(),
		},
	}
	lambda.Start(l.handler)
}
