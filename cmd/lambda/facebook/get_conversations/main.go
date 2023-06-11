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

var errNoPageIDPath = errors.New("err path parameter parameters page_id not given")

func (c *config) handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	discord.Log(c.DiscordWebhookURL, "facebook get conversations handler")

	ctx, cancel := context.WithTimeout(context.Background(), 2500*time.Millisecond)
	defer cancel()

	pathParams := request.PathParameters
	// shopID := pathParams["shop_id"]
	pageID, ok := pathParams["page_id"]
	if !ok {
		discord.Log(c.DiscordWebhookURL, "err path param page_id was not given")
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Bad Request",
		}, errNoPageIDPath
	}

	stdConversations, err := c.DbClient.QueryConversations(ctx, pageID)
	if err != nil {
		discord.Log(c.DiscordWebhookURL, fmt.Sprint(err))
		return events.APIGatewayProxyResponse{
			StatusCode: 502,
			Body:       "Bad Gateway",
		}, err
	}

	jsonBodyByte, err := json.Marshal(stdConversations)
	if err != nil {
		discord.Log(c.DiscordWebhookURL, fmt.Sprint(err))
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
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
		DiscordWebhookURL: os.Getenv("DISCORD_WEBHOOK_URL"),
		DbClient:          dbClient,
	}
	defer func() {
		discord.Log(c.DiscordWebhookURL, "defer dbclient close")
		c.DbClient.Close(ctx)
	}()

	lambda.Start(c.handler)
}
