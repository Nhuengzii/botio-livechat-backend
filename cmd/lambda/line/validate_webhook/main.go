package main

import (
	"context"
	"errors"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/sqswrapper"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"os"
)

func (c *config) handler(ctx context.Context, req events.APIGatewayProxyRequest) (_ events.APIGatewayProxyResponse, err error) {
	defer func() {
		if err != nil {
			logMessage := "lambda/line/validate_webhook/main.config.handler: " + err.Error()
			log.Println(logMessage)
			discord.Log(c.discordWebhookURL, logMessage)
		}
	}()
	//pathParameters := req.PathParameters
	//shopID := pathParameters["shop_id"]
	//pageID := pathParameters["page_id"]
	//lineChannelSecret := // TODO get from db
	lineSignature := req.Headers["x-line-signature"]
	webhookBodyString := req.Body
	err = validateSignature(c.lineChannelSecret, lineSignature, webhookBodyString)
	if err != nil {
		if errors.Is(err, errInvalidSignature) {
			return events.APIGatewayProxyResponse{
				StatusCode: 401,
				Body:       "Unauthorized",
			}, err
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, err
	}
	err = c.sqsClient.SendMessage(c.sqsQueueURL, webhookBodyString)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "OK",
	}, nil
}

func main() {
	c := &config{
		discordWebhookURL: os.Getenv("DISCORD_WEBHOOK_URL"),
		sqsQueueURL:       os.Getenv("SQS_QUEUE_URL"),
		lineChannelSecret: os.Getenv("LINE_CHANNEL_SECRET"), // TODO remove and get from some db with shopID and pageID
		sqsClient:         sqswrapper.NewClient(os.Getenv("AWS_REGION")),
	}
	lambda.Start(c.handler)
}
