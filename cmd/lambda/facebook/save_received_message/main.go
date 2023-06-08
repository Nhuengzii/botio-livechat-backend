package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/Nhuengzii/botio-livechat-backend/internal/db"
	"github.com/Nhuengzii/botio-livechat-backend/internal/discord"
	"github.com/Nhuengzii/botio-livechat-backend/pkg/stdmessage"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/bson"
)

type Lambda struct {
	config
}

type recievedSqsMessage struct {
	Message string `json:"Message"`
}

var (
	errUnmarshalRecievedBody    = errors.New("Error json unmarshal recieve body")
	errUnmarshalRecievedMessage = errors.New("Error json unmarshal recieve message")
)

func (l Lambda) handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	var recieveBody recievedSqsMessage
	var recieveMessage stdmessage.StdMessage
	for _, record := range sqsEvent.Records {
		err := json.Unmarshal([]byte(record.Body), &recieveBody)
		if err != nil {
			discord.Log(l.DiscordWebhookURL, "Error unmarshal recieveBody")
			return errUnmarshalRecievedBody
		}
		err = bson.UnmarshalExtJSON([]byte(recieveBody.Message), true, &recieveMessage)
		if err != nil {
			return errUnmarshalRecievedMessage
		}

		// implement update conversation
		// implement add message
	}
	return nil
}

func main() {
	dbClient, err := db.NewClient(context.TODO(), &db.Target{
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
			DiscordWebhookURL: os.Getenv("DISCORD_WEBHOOK_URL"),
			DbClient:          dbClient,
		},
	}
	lambda.Start(l.handler)
}
