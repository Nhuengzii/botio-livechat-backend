package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/Nhuengzii/botio-livechat-backend/livechat"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/bson"
)

type Lambda struct {
	config
}

type receivedMessage struct {
	Message string `json:"Message"`
}

var (
	errUnmarshalReceivedBody    = errors.New("Error json unmarshal recieve body")
	errUnmarshalReceivedMessage = errors.New("Error json unmarshal recieve message")
)

func (l Lambda) handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	var receiveBody receivedMessage
	var receiveMessage livechat.StdMessage
	for _, record := range sqsEvent.Records {
		err := json.Unmarshal([]byte(record.Body), &receiveBody)
		if err != nil {
			discord.Log(l.DiscordWebhookURL, "Error unmarshal receiveBody")
			return errUnmarshalReceivedBody
		}
		err = bson.UnmarshalExtJSON([]byte(receiveBody.Message), true, &receiveMessage)
		if err != nil {
			discord.Log(l.DiscordWebhookURL, "Error unmarshal receiveMessage")
			return errUnmarshalReceivedMessage
		}
		err = l.config.DbClient.UpdateConversationOnNewMessage(ctx, &receiveMessage)
		if err != nil {
			if errors.Is(err, mongodb.ErrNoConversations) {
				conversation, err := newStdConversation(l.config.FacebookAccessToken, &receiveMessage)
				if err != nil {
					return err
				}
				err = l.config.DbClient.InsertConversation(ctx, conversation)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
		err = l.config.DbClient.InsertMessage(ctx, &receiveMessage)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func main() {
	dbClient, err := mongodb.NewClient(context.TODO(), &mongodb.Target{
		URI:                     os.Getenv("DATABASE_CONNECTION_URI"),
		Database:                "BotioLivechat",
		CollectionMessages:      "facebook_messages",
		CollectionConversations: "facebook_conversations",
	})
	if err != nil {
		return
	}
	l := Lambda{
		config: config{
			DiscordWebhookURL:   os.Getenv("DISCORD_WEBHOOK_URL"),
			DbClient:            dbClient,
			FacebookAccessToken: os.Getenv("ACCESS_TOKEN"),
		},
	}
	defer func() {
		discord.Log(l.DiscordWebhookURL, "defer dbclient close")
		l.DbClient.Close(context.TODO())
	}()

	lambda.Start(l.handler)
}
