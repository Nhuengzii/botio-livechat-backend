package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/snswrapper"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/storage/amazons3"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/v7/linebot"
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
		page, err := c.dbClient.GetLineAuthentication(ctx, pageID)
		if err != nil {
			return err
		}
		lineChannelSecret := page.Secret
		lineChannelAccessToken := page.AccessToken
		bot, err := linebot.New(lineChannelSecret, lineChannelAccessToken)
		if err != nil {
			return err
		}
		shop, err := c.dbClient.GetShop(ctx, pageID)
		if err != nil {
			return err
		}
		shopID := shop.ShopID
		err = c.handleEvents(ctx, shopID, pageID, bot, hookBody)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*2500)
	defer cancel()
	var (
		mongodbURI        = os.Getenv("MONGODB_URI")
		mongodbDatabase   = os.Getenv("MONGODB_DATABASE")
		discordWebhookURL = os.Getenv("DISCORD_WEBHOOK_URL")
		snsTopicARN       = os.Getenv("SNS_TOPIC_ARN")
		awsRegion         = os.Getenv("AWS_REGION")
		s3BucketName      = os.Getenv("S3_BUCKET_NAME")
	)
	dbClient, err := mongodb.NewClient(ctx, mongodb.Target{
		URI:                     mongodbURI,
		Database:                mongodbDatabase,
		CollectionConversations: "conversations",
		CollectionMessages:      "messages",
		CollectionShops:         "shops",
		CollectionTemplates:     "templates",
	})
	if err != nil {
		logMessage := "cmd/lambda/line/standardize_webhook/main.main: " + err.Error()
		discord.Log(discordWebhookURL, logMessage)
		log.Fatalln(logMessage)
	}

	stroageClient := amazons3.NewClient(awsRegion, s3BucketName)
	defer dbClient.Close(ctx)
	c := &config{
		discordWebhookURL: discordWebhookURL,
		snsTopicARN:       snsTopicARN,
		snsClient:         snswrapper.NewClient(awsRegion),
		dbClient:          dbClient,
		storageClient:     stroageClient,
	}
	lambda.Start(c.handler)
}
