package main

import (
	"context"
	"log"
	"os"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func (c *config) handler(ctx context.Context, req events.APIGatewayProxyRequest) (_ events.APIGatewayProxyResponse, err error) {
	defer func() {
		if err != nil {
			logMessage := "lambda/line/post_message: " + err.Error()
			log.Println(logMessage)
			discord.Log(c.discordWebhookURL, logMessage)
		}
	}()
	// pathParameters := req.PathParameters
	// pageID := pathParameters["page_id"]
	// conversationID := pathParameters["conversation_id"]
	// bot, err := linebot.New(c.lineChannelSecret, c.lineChannelAccessToken)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "OK",
	}, nil
}

func main() {
	c := &config{
		discordWebhookURL:      os.Getenv("DISCORD_WEBHOOK_URL"),
		lineChannelSecret:      os.Getenv("LINE_CHANNEL_SECRET"),
		lineChannelAccessToken: os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
		dbClient:               nil,
	}
	lambda.Start(c.handler)
}
