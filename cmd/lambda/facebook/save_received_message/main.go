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
			discord.Log(c.DiscordWebhookURL, fmt.Sprint(err))
		}
	}()

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
		err = c.DbClient.UpdateConversationOnNewMessage(ctx, &receiveMessage)
		if err != nil {
			if errors.Is(err, mongodb.ErrNoDocuments) {
				conversation, err := c.newStdConversation(ctx, &receiveMessage)
				if err != nil {
					return err
				}
				err = c.DbClient.InsertConversation(ctx, conversation)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
		err = c.DbClient.InsertMessage(ctx, &receiveMessage)
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
	})
	if err != nil {
		return
	}
	c := config{
		DiscordWebhookURL:   os.Getenv("DISCORD_WEBHOOK_URL"),
		DbClient:            dbClient,
		FacebookAccessToken: os.Getenv("ACCESS_TOKEN"),
	}

	defer func() {
		discord.Log(c.DiscordWebhookURL, "defer dbclient close")
		c.DbClient.Close(ctx)
	}()

	lambda.Start(c.handler)
}
