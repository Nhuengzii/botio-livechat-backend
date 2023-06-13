package main

import (
	"context"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
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
			logMessage := "cmd/lambda/line/standardize_webhook/main.config.handler: " + err.Error()
			log.Println(logMessage)
			discord.Log(c.discordWebhookURL, logMessage)
		}
	}()
	for _, sqsMessage := range sqsEvent.Records {
		hookBody, err := parseWebhookBody(sqsMessage.Body)
		if err != nil {
			return err
		}
		pageID := hookBody.Destination
		shop, err := c.dbClient.QueryShop(ctx, pageID)
		shopID := shop.ShopID
		err = c.handleEvents(ctx, shopID, pageID, hookBody)
		if err != nil {
			return err
		}
	}
	return nil
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
		log.Fatalln("cmd/lambda/line/standardize_webhook/main.main: " + err.Error())
	}
	defer dbClient.Close(ctx)
	c := &config{
		discordWebhookURL: os.Getenv("DISCORD_WEBHOOK_URL"),
		snsTopicARN:       os.Getenv("SNS_TOPIC_ARN"),
		snsClient:         snswrapper.NewClient(os.Getenv("AWS_REGION")),
		dbClient:          dbClient,
	}
	lambda.Start(c.handler)
}
