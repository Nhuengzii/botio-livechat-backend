package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
	"os"
	"time"

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

func (c *config) handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	var receiveBody receivedMessage
	var receiveMessage stdmessage.StdMessage
	for _, record := range sqsEvent.Records {
		err := json.Unmarshal([]byte(record.Body), &receiveBody)
		if err != nil {
			discord.Log(c.DiscordWebhookURL, "Error unmarshal receiveBody")
			return errUnmarshalReceivedBody
		}
		err = bson.UnmarshalExtJSON([]byte(receiveBody.Message), true, &receiveMessage)
		if err != nil {
			discord.Log(c.DiscordWebhookURL, "Error unmarshal receiveMessage")
			return errUnmarshalReceivedMessage
		}
		err = c.DbClient.UpdateConversationOnNewMessage(ctx, &receiveMessage)
		if err != nil {
			if errors.Is(err, mongodb.ErrNoConversations) {
				conversation, err := newStdConversation(c.FacebookAccessToken, &receiveMessage)
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
		URI:                     os.Getenv("DATABASE_CONNECTION_URI"),
		Database:                "BotioLivechat",
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
