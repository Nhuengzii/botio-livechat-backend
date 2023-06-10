package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func (c *config) handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	discord.Log(c.DiscordWebhookURL, "facebook get messages handler")

	ctx, cancel := context.WithTimeout(context.Background(), 2500*time.Millisecond)
	defer cancel()

	pathParams := request.PathParameters
	// shopID := pathParams["shop_id"]
	pageID := pathParams["page_id"]
	conversationID := pathParams["conversation_id"]

	stdMessages, err := c.DbClient.QueryMessages(ctx, pageID, conversationID)
	if err != nil {
		discord.Log(c.DiscordWebhookURL, fmt.Sprint(err))
		return events.APIGatewayProxyResponse{
			StatusCode: 502,
			Body:       "Bad Gateway",
		}, err
	}

	jsonBodyByte, err := json.Marshal(*stdMessages)
	if err != nil {
		discord.Log(c.DiscordWebhookURL, fmt.Sprint(err))
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, err
	}

	err = c.DbClient.UpdateConversationIsRead(ctx, conversationID)
	if err != nil {
		discord.Log(c.DiscordWebhookURL, fmt.Sprint(err))
		return events.APIGatewayProxyResponse{
			StatusCode: 502,
			Body:       "Bad Gateway",
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonBodyByte),
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
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
	c := config{
		DiscordWebhookURL: os.Getenv("DISCORD_WEBHOOK_URL"),
		DbClient:          dbClient,
	}

	lambda.Start(c.handler)
}
