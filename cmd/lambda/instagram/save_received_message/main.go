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
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type receivedMessage struct {
	Message string `json:"Message"`
}

var (
	errUnmarshalReceivedBody    = errors.New("error json unmarshal receive body")
	errUnmarshalReceivedMessage = errors.New("error json unmarshal receive message")
	errUnsupportedUserType      = errors.New("error unsupported user type")
	errCannotGetUserPSID        = errors.New("error cannot get ig user's PSID")
)

func (c *config) handler(ctx context.Context, sqsEvent events.SQSEvent) (err error) {
	defer func() {
		if err != nil {
			discord.Log(c.discordWebhookURL, fmt.Sprint(err))
		}
	}()

	var receiveBody receivedMessage
	var receiveMessage stdmessage.StdMessage
	for _, record := range sqsEvent.Records {
		err := json.Unmarshal([]byte(record.Body), &receiveBody)
		if err != nil {
			return errUnmarshalReceivedBody
		}
		err = json.Unmarshal([]byte(receiveBody.Message), &receiveMessage)
		if err != nil {
			return errUnmarshalReceivedMessage
		}

		if !receiveMessage.IsDeleted {
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
		} else {
			err = c.dbClient.UpdateConversationOnDeletedMessage(ctx, receiveMessage)
			if err != nil {
				return err
			}
			err = c.dbClient.RemoveDeletedMessage(ctx, receiveMessage.ShopID, stdmessage.PlatformInstagram, receiveMessage.ConversationID, receiveMessage.MessageID)
			if err != nil {
				return err
			}
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
		CollectionMessages:      "messages",
		CollectionConversations: "conversations",
		CollectionShops:         "shops",
		CollectionTemplates:     "templates",
	})
	c := config{
		discordWebhookURL: discordWebhookURL,
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
