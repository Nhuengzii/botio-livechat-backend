package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/snswrapper"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	errNoMessageEntry       = errors.New("error! no message entry")
	errUnknownWebhookType   = errors.New("error! unknown webhook type found")
	errUnknownWebhookObject = errors.New("error! unknown webhook Object found")
	errUnknownTemplateType  = errors.New("error! unknown attachment template type")
)

func (c *config) handler(ctx context.Context, sqsEvent events.SQSEvent) (err error) {
	defer func() {
		if err != nil {
			discord.Log(c.discordWebhookURL, fmt.Sprint(err))
		}
	}()

	var receiveWebhook ReceiveWebhook

	for _, record := range sqsEvent.Records {
		err := json.Unmarshal([]byte(record.Body), &receiveWebhook)
		if err != nil {
			return err
		}
		err = c.handleReceiveWebhook(ctx, &receiveWebhook)
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
	)

	dbClient, err := mongodb.NewClient(ctx, mongodb.Target{
		URI:                     mongodbURI,
		Database:                mongodbDatabase,
		CollectionMessages:      "messages",
		CollectionConversations: "conversations",
		CollectionShops:         "shops",
	})
	c := config{
		discordWebhookURL: discordWebhookURL,
		snsTopicARN:       snsTopicARN,
		snsClient:         snswrapper.NewClient(awsRegion),
		dbClient:          dbClient,
	}
	if err != nil {
		discord.Log(c.discordWebhookURL, fmt.Sprintln(err))
		log.Fatalln(err)
	}
	defer func() {
		discord.Log(c.discordWebhookURL, "defer dbClient close")
		c.dbClient.Close(ctx)
	}()
	lambda.Start(c.handler)
}
