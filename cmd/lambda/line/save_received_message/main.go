package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/snswrapper"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func (c *config) handler(ctx context.Context, sqsEvent events.SQSEvent) (err error) {
	defer func() {
		if err != nil {
			logMessage := "cmd/lambda/line/save_received_message/main.config.handler: " + err.Error()
			log.Println(logMessage)
			discord.Log(c.discordWebhookURL, logMessage)
		}
	}()
	for _, sqsMessage := range sqsEvent.Records {
		snsMessageString := sqsMessage.Body
		var snsMessage snswrapper.SNSMessage
		err = json.Unmarshal([]byte(snsMessageString), &snsMessage)
		if err != nil {
			return err
		}
		var stdMessage stdmessage.StdMessage
		err = json.Unmarshal([]byte(snsMessage.Message), &stdMessage)
		if err != nil {
			return err
		}
		pageID := stdMessage.PageID
		shop, err := c.dbClient.QueryLinePage(ctx, pageID)
		if err != nil {
			return err
		}
		lineChannelAccessToken := shop.AccessToken
		err = c.handleMessage(ctx, lineChannelAccessToken, &stdMessage)
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
	)
	dbClient, err := mongodb.NewClient(ctx, mongodb.Target{
		URI:                     mongodbURI,
		Database:                mongodbDatabase,
		CollectionConversations: "conversations",
		CollectionMessages:      "messages",
		CollectionShops:         "shops",
	})
	if err != nil {
		logMessage := "cmd/lambda/line/save_received_message/main.main: " + err.Error()
		discord.Log(discordWebhookURL, logMessage)
		log.Fatalln(logMessage)
	}
	defer dbClient.Close(ctx)
	c := &config{
		discordWebhookURL: discordWebhookURL,
		dbClient:          dbClient,
	}
	lambda.Start(c.handler)
}
