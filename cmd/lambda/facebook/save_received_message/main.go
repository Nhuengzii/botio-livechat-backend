package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/bson"
)

type receivedMessage struct {
	Message string `json:"Message"`
}

var (
	errUnmarshalReceivedBody    = errors.New("Error json unmarshal recieve body")
	errUnmarshalReceivedMessage = errors.New("Error json unmarshal recieve message")
)

func (c *config) handler(ctx context.Context, sqsEvent events.SQSEvent) (err error) {
	defer func() {
		if err != nil {
			discord.Log(c.discordWebhookUrl, fmt.Sprint(err))
		}
	}()

	discord.Log(c.discordWebhookUrl, "facebook save recieved message handler")

	var receiveBody receivedMessage
	var receiveMessage stdmessage.StdMessage
	for _, record := range sqsEvent.Records {
		err := json.Unmarshal([]byte(record.Body), &receiveBody)
		if err != nil {
			return errUnmarshalReceivedBody
		}
		err = bson.UnmarshalExtJSON([]byte(receiveBody.Message), true, &receiveMessage)
		if err != nil {
			return errUnmarshalReceivedMessage
		}
		err = c.dbClient.UpdateConversationOnNewMessage(ctx, &receiveMessage)
		if err != nil {
			if errors.Is(err, mongodb.ErrNoDocuments) {
				conversation, err := c.newStdConversation(ctx, &receiveMessage)
				if err != nil {
					return err
				}
				err = c.dbClient.InsertConversation(ctx, conversation)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
		err = c.dbClient.InsertMessage(ctx, &receiveMessage)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*2500)
	defer cancel()
	dbClient, err := mongodb.NewClient(ctx, &mongodb.Target{
		URI:                     os.Getenv("MONGODB_URI"),
		Database:                os.Getenv("MONGODB_DATABASE"),
		CollectionMessages:      "facebook_messages",
		CollectionConversations: "facebook_conversations",
		CollectionShops:         "shops",
	})
	if err != nil {
		return
	}
	c := config{
		discordWebhookUrl:   os.Getenv("DISCORD_WEBHOOK_URL"),
		dbClient:            dbClient,
		facebookAccessToken: os.Getenv("ACCESS_TOKEN"),
	}

	defer func() {
		discord.Log(c.discordWebhookUrl, "defer dbclient close")
		c.dbClient.Close(ctx)
	}()

	lambda.Start(c.handler)
}
