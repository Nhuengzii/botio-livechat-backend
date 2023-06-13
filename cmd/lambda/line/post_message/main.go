package main

import (
	"context"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
	"log"
	"os"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func (c *config) handler(ctx context.Context, req events.APIGatewayProxyRequest) (_ events.APIGatewayProxyResponse, err error) {
	defer func() {
		if err != nil {
			logMessage := "cmd/lambda/line/post_message.main.config.handler: " + err.Error()
			log.Println(logMessage)
			discord.Log(c.discordWebhookURL, logMessage)
		}
	}()
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
		log.Fatalln("cmd/lambda/line/post_message.main.main: " + err.Error())
	}
	defer dbClient.Close(ctx)
	c := &config{
		discordWebhookURL: os.Getenv("DISCORD_WEBHOOK_URL"),
		dbClient:          dbClient,
	}
	lambda.Start(c.handler)
}
