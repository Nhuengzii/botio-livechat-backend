package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	errNoPageIDPath         = errors.New("err path parameter parameters page_id not given")
	errNoConversationIDPath = errors.New("err path parameter conversation_id not given")
)

func (c *config) handler(ctx context.Context, request events.APIGatewayProxyRequest) (_ events.APIGatewayProxyResponse, err error) {
	discord.Log(c.discordWebhookURL, "facebook get messages handler")

	defer func() {
		if err != nil {
			discord.Log(c.discordWebhookURL, fmt.Sprintln(err))
		}
	}()

	pathParams := request.PathParameters
	// shopID := pathParams["shop_id"]
	pageID, ok := pathParams["page_id"]
	if !ok {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Bad Request",
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, errNoPageIDPath
	}
	conversationID, ok := pathParams["conversation_id"]
	if !ok {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Bad Request",
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, errNoConversationIDPath
	}

	stdMessages, err := c.dbClient.QueryMessages(ctx, pageID, conversationID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 502,
			Body:       "Bad Gateway",
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, err
	}

	jsonBodyByte, err := json.Marshal(stdMessages)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, err
	}

	err = c.dbClient.UpdateConversationIsRead(ctx, conversationID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 502,
			Body:       "Bad Gateway",
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
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
	ctx, cancel := context.WithTimeout(context.Background(), 2500*time.Millisecond)
	defer cancel()
	dbClient, err := mongodb.NewClient(ctx, mongodb.Target{
		URI:                     os.Getenv("MONGODB_URI"),
		Database:                os.Getenv("MONGODB_DATABASE"),
		CollectionMessages:      "messages",
		CollectionConversations: "conversations",
	})
	if err != nil {
		return
	}
	c := config{
		discordWebhookURL: os.Getenv("DISCORD_WEBHOOK_URL"),
		dbClient:          dbClient,
	}
	defer func() {
		discord.Log(c.discordWebhookURL, "defer dbClient close")
		c.dbClient.Close(ctx)
	}()

	lambda.Start(c.handler)
}
