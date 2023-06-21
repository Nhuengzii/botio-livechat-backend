package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/apigateway"
	"log"
	"os"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/getconversation"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func (c *config) handler(ctx context.Context, req events.APIGatewayProxyRequest) (_ events.APIGatewayProxyResponse, err error) {
	defer func() {
		if err != nil {
			logMessage := "cmd/lambda/line/get_conversation/main.config.handler: " + err.Error()
			log.Println(logMessage)
			discord.Log(c.discordWebhookURL, logMessage)
		}
	}()
	pathParameters := req.PathParameters
	shopID := pathParameters["shop_id"]
	pageID := pathParameters["page_id"]
	conversationID := pathParameters["conversation_id"]
	conversation, err := c.dbClient.QueryConversation(ctx, shopID, pageID, conversationID)
	if err != nil {
		if errors.Is(err, mongodb.ErrNoDocuments) {
			emptyResponse := getconversation.Response{
				Conversation: nil,
			}
			emptyResponseJSON, err := json.Marshal(emptyResponse)
			if err != nil {
				return apigateway.NewProxyResponse(500, "Internal ServerError", "*"), err
			}
			return apigateway.NewProxyResponse(200, string(emptyResponseJSON), "*"), err
		}
		return apigateway.NewProxyResponse(500, "Internal ServerError", "*"), err
	}
	resp := getconversation.Response{
		Conversation: conversation,
	}
	responseJSON, err := json.Marshal(resp)
	if err != nil {
		return apigateway.NewProxyResponse(500, "Internal ServerError", "*"), err
	}
	return apigateway.NewProxyResponse(200, string(responseJSON), "*"), nil
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
		logMessage := "cmd/lambda/line/get_conversations/main.main: " + err.Error()
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
