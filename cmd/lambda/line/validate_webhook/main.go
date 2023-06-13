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
	"time"
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
	if err != nil {
		if errors.Is(err, mongodb.ErrNoDocuments) {
			return events.APIGatewayProxyResponse{
				StatusCode: 404,
				Headers: map[string]string{
					"Access-Control-Allow-Origin": "*",
				},
				Body: "Not Found",
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*2500)
	defer cancel()
	var (
		mongodbURI        = os.Getenv("MONGODB_URI")
		mongodbDatabase   = os.Getenv("MONGODB_DATABASE")
		discordWebhookURL = os.Getenv("DISCORD_WEBHOOK_URL")
		sqsQueueURL       = os.Getenv("SQS_QUEUE_URL")
		awsRegion         = os.Getenv("AWS_REGION")
	)
	dbClient, err := mongodb.NewClient(ctx, mongodb.Target{
		URI:                     mongodbURI,
		Database:                mongodbDatabase,
		CollectionConversations: "conversations",
		CollectionMessages:      "messages",
		CollectionShops:         "shops",
	})
	if err != nil {
		logMessage := "cmd/lambda/line/validate_webhook/main.main: " + err.Error()
		discord.Log(os.Getenv(discordWebhookURL), logMessage)
		log.Fatalln(logMessage)
	}
	defer dbClient.Close(ctx)
	c := &config{
		discordWebhookURL: discordWebhookURL,
		sqsQueueURL:       os.Getenv(sqsQueueURL),
		sqsClient:         sqswrapper.NewClient(awsRegion),
		dbClient:          dbClient,
	}
	lambda.Start(c.handler)
}
