package main

import (
	"context"
	"encoding/json"
	"github.com/Nhuengzii/botio-livechat-backend/livechat"
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
			logMessage := "lambda/line/save_received_message/main.config.handler: " + err.Error()
			log.Println(logMessage)
			discord.Log(c.discordWebhookURL, logMessage)
		}
	}()
	if c.dbClient == nil {
		c.dbClient, err = mongodb.NewClient(ctx, &mongodb.Target{
			URI:                     c.mongodbURI,
			Database:                c.mongodbDatabase,
			CollectionConversations: c.mongodbCollectionLineConversations,
			CollectionMessages:      c.mongodbCollectionLineMessages,
		})
		if err != nil {
			return err
		}
	}
	for _, sqsMessage := range sqsEvent.Records {
		snsMessageString := sqsMessage.Body
		var snsMessage *snswrapper.SNSMessage
		err = json.Unmarshal([]byte(snsMessageString), snsMessage)
		if err != nil {
			return err
		}
		var stdMessage *livechat.StdMessage
		err = json.Unmarshal([]byte(snsMessage.Message), stdMessage)
		if err != nil {
			return err
		}
		// TODO get lineChannelAccessToken from db with shopID and pageID here and pass to updateDB through a parameter
		err = updateDB(ctx, c, stdMessage)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	c := &config{
		discordWebhookURL:                  os.Getenv("DISCORD_WEBHOOK_URL"),
		lineChannelAccessToken:             os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"), // TODO remove and get from db with shopID and pageID
		mongodbURI:                         os.Getenv("MONGODB_URI"),
		mongodbDatabase:                    os.Getenv("MONGODB_DATABASE"),
		mongodbCollectionLineConversations: os.Getenv("MONGODB_COLLECTION_LINE_CONVERSATIONS"),
		mongodbCollectionLineMessages:      os.Getenv("MONGODB_COLLECTION_LINE_MESSAGES"),
		dbClient:                           nil,
	}
	lambda.Start(c.handler)
}
