package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/getall"
	"log"
	"os"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/apigateway"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func (c *config) handler(ctx context.Context, req events.APIGatewayProxyRequest) (_ events.APIGatewayProxyResponse, err error) {
	defer func() {
		if err != nil {
			logMessage := "cmd/lambda/all/get_conversations/main.config.handler: " + err.Error()
			log.Println(logMessage)
			discord.Log(c.discordWebhookURL, logMessage)
		}
	}()

	pathParameters := req.PathParameters
	shopID := pathParameters["shop_id"]

	platformStatuses, err := c.dbClient.ListShopPlatformsStatuses(ctx, shopID)
	if err != nil {
		if errors.Is(err, mongodb.ErrNoDocuments) {
			return apigateway.NewProxyResponse(404, "Not Found", "*"), nil
		}
		return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
	}
	if len(platformStatuses) == 0 {
		return apigateway.NewProxyResponse(204, "No Content", "*"), nil
	}

	response := getall.Response{
		Statuses: platformStatuses,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
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
		logMessage := "cmd/lambda/all/get_conversations/main.main: " + err.Error()
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
