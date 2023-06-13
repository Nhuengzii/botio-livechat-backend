package main

import (
	"context"
	"errors"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
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
			logMessage := "cmd/lambda/line/validate_webhook/main.config.handler: " + err.Error()
			log.Println(logMessage)
			discord.Log(c.discordWebhookURL, logMessage)
		}
	}()
	pathParameters := req.PathParameters
	pageID := pathParameters["page_id"]
	shop, err := c.dbClient.QueryLinePage(ctx, pageID)
	lineChannelSecret := shop.Secret
	lineSignature := req.Headers["x-line-signature"]
	webhookBodyString := req.Body
	err = validateSignature(lineChannelSecret, lineSignature, webhookBodyString)
	if err != nil {
		if errors.Is(err, errInvalidSignature) {
			return events.APIGatewayProxyResponse{
				StatusCode: 401,
				Headers: map[string]string{
					"Access-Control-Allow-Origin": "*",
				},
				Body: "Unauthorized",
			}, err
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
			Body: "Internal Server Error",
		}, err
	}
	err = c.sqsClient.SendMessage(c.sqsQueueURL, webhookBodyString)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
			Body: "Internal Server Error",
		}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
		Body: "OK",
	}, nil
}

func main() {
	ctx := context.Background()
	dbClient, err := mongodb.NewClient(ctx, &mongodb.Target{
		URI:                     os.Getenv("MONGODB_URI"),
		Database:                os.Getenv("MONGODB_DATABASE"),
		CollectionConversations: "conversations",
		CollectionMessages:      "messages",
		CollectionShops:         "shops",
	})
	if err != nil {
		log.Fatalln("cmd/lambda/line/validate_webhook/main.main: " + err.Error())
	}
	defer dbClient.Close(ctx)
	c := &config{
		discordWebhookURL: os.Getenv("DISCORD_WEBHOOK_URL"),
		sqsQueueURL:       os.Getenv("SQS_QUEUE_URL"),
		sqsClient:         sqswrapper.NewClient(os.Getenv("AWS_REGION")),
		dbClient:          dbClient,
	}
	lambda.Start(c.handler)
}
