package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

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
	ctx := context.Background()
	dbClient, err := mongodb.NewClient(ctx, mongodb.Target{
		URI:                     os.Getenv("MONGODB_URI"),
		Database:                os.Getenv("MONGODB_DATABASE"),
		CollectionConversations: "conversations",
		CollectionMessages:      "messages",
		CollectionShops:         "shops",
	})
	if err != nil {
		log.Fatalln("cmd/lambda/line/save_received_message/main.main: " + err.Error())
	}
	defer dbClient.Close(ctx)
	c := &config{
		discordWebhookURL: os.Getenv("DISCORD_WEBHOOK_URL"),
		dbClient:          dbClient,
	}
	lambda.Start(c.handler)
}
