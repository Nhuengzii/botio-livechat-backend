package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/getmessages"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func (c *config) handler(ctx context.Context, req events.APIGatewayProxyRequest) (_ events.APIGatewayProxyResponse, err error) {
	defer func() {
		if err != nil {
			logMessage := "cmd/lambda/line/get_messages/main.config.handler: " + err.Error()
			log.Println(logMessage)
			discord.Log(c.discordWebhookURL, logMessage)
		}
	}()
	pathParameters := req.PathParameters
	shopID := pathParameters["shop_id"]
	pageID := pathParameters["page_id"]
	conversationID := pathParameters["conversation_id"]
	messages, err := c.dbClient.QueryMessages(ctx, shopID, pageID, conversationID)
	if err != nil {
		if errors.Is(err, mongodb.ErrNoDocuments) {
			return events.APIGatewayProxyResponse{
				StatusCode: 404,
				Headers: map[string]string{
					"Access-Control-Allow-Origin": "*",
				},
				Body: "Not Found",
			}, err
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
			Body: "Internal Server Error",
		}, err
	}
	if len(messages) != 0 {
		err = c.dbClient.UpdateConversationIsRead(ctx, conversationID) // TODO remove this UpdateConversationIsRead and do with another api endpoint
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Headers: map[string]string{
					"Access-Control-Allow-Origin": "*",
				},
				Body: "Internal Server Error",
			}, err
		}
	}
	resp := getmessages.Response{
		Messages: messages,
	}
	responseJSON, err := json.Marshal(resp)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
			Body: "Internal Server Error",
		}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
		Body: string(responseJSON),
	}, nil
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
		CollectionConversations: "conversations",
		CollectionMessages:      "messages",
		CollectionShops:         "shops",
	})
	if err != nil {
		logMessage := "cmd/lambda/line/get_messages/main.main: " + err.Error()
		discord.Log(discordWebhookURL, logMessage)
		log.Fatalln(logMessage)
	}
	defer dbClient.Close(ctx)
	c := &config{
		discordWebhookURL: discordWebhookURL,
		dbClient:          dbClient,
	}
	lambda.Start(c.handler)
}
