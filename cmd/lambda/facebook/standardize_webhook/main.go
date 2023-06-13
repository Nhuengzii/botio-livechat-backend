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
	errNoMessageEntry       = errors.New("Error! no message entry")
	errUnknownWebhookType   = errors.New("Error! unknown webhook type found!")
	errUnknownWebhookObject = errors.New("Error! unknown webhook Object found!")
)

func (c *config) handler(ctx context.Context, sqsEvent events.SQSEvent) (err error) {
	defer func() {
		if err != nil {
			discord.Log(c.discordWebhookURL, fmt.Sprint(err))
		}
	}()

	discord.Log(c.discordWebhookURL, "facebook standardize webhook handler")
	start := time.Now()
	var recieveWebhook ReceiveWebhook

	for _, record := range sqsEvent.Records {
		err := json.Unmarshal([]byte(record.Body), &recieveWebhook)
		if err != nil {
			return errUnknownWebhookObject
		}
		discord.Log(c.discordWebhookURL, fmt.Sprintf("%+v\n", recieveWebhook))
		err = c.handleRecieveWebhook(ctx, &recieveWebhook)
		if err != nil {
			return err
		}
	}
	discord.Log(c.discordWebhookURL, fmt.Sprintf("Elapsed: %v", time.Since(start)))
	return nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*2500)
	defer cancel()

	dbClient, err := mongodb.NewClient(ctx, mongodb.Target{
		URI:                     os.Getenv("MONGODB_URI"),
		Database:                os.Getenv("MONGODB_DATABASE"),
		CollectionMessages:      "facebook_messages",
		CollectionConversations: "facebook_conversations",
		CollectionShops:         "shops",
	})
	if err != nil {
		log.Println(err)
		return
	}
	c := config{
		discordWebhookURL: os.Getenv("DISCORD_WEBHOOK_URL"),
		snsTopicARN:       os.Getenv("SNS_TOPIC_ARN"),
		snsClient:         snswrapper.NewClient(os.Getenv("AWS_REGION")),
		dbClient:          dbClient,
	}
	defer func() {
		discord.Log(c.discordWebhookURL, "defer dbclient close")
		c.dbClient.Close(ctx)
	}()
	lambda.Start(c.handler)
}
