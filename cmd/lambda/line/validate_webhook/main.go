package main

import (
	"context"
	"errors"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/apigateway"
	"log"
	"os"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/sqswrapper"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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
	shop, err := c.dbClient.GetLineAuthentication(ctx, pageID)
	if err != nil {
		if errors.Is(err, mongodb.ErrNoDocuments) {
			return apigateway.NewProxyResponse(404, "Not Found", "*"), err
		}
		return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
	}
	lineChannelSecret := shop.Secret
	lineSignature := req.Headers["x-line-signature"]
	webhookBodyString := req.Body
	err = validateSignature(lineChannelSecret, lineSignature, webhookBodyString)
	if err != nil {
		if errors.Is(err, errInvalidSignature) {
			return apigateway.NewProxyResponse(401, "Unauthorized", "*"), err
		}
		return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
	}
	err = c.sqsClient.SendMessage(c.sqsQueueURL, webhookBodyString)
	if err != nil {
		return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
	}
	return apigateway.NewProxyResponse(200, "OK", "*"), nil
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
		sqsQueueURL:       sqsQueueURL,
		sqsClient:         sqswrapper.NewClient(awsRegion),
		dbClient:          dbClient,
	}
	lambda.Start(c.handler)
}
