package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/shops"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/apigateway"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func (c *config) handler(ctx context.Context, req events.APIGatewayProxyRequest) (_ events.APIGatewayProxyResponse, err error) {
	defer func() {
		if err != nil {
			logMessage := "cmd/lambda/shops/post_shops/main.config.handler: " + err.Error()
			log.Println(logMessage)
			discord.Log(c.discordWebhookURL, logMessage)
		}
	}()
	reqBody := req.Body
	if reqBody == "" {
		return apigateway.NewProxyResponse(400, "Bad Request: Request body must not be empty", "*"), nil
	}

	shop := shops.Shop{}
	err = json.Unmarshal([]byte(reqBody), &shop)
	if err != nil {
		return apigateway.NewProxyResponse(400, "Bad : Check request body format", "*"), nil
	}

	err = c.dbClient.CheckShopExists(ctx, shop.ShopID)
	if err == nil {
		return apigateway.NewProxyResponse(400, "Bad Request: Shop already exists", "*"), nil
	} else {
		if err != mongodb.ErrNoDocuments {
			return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
		}
	}

	err = c.dbClient.InsertShop(ctx, shop)
	if err != nil {
		return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
	}

	return apigateway.NewProxyResponse(200, "OK", "*"), nil
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
		logMessage := "cmd/lambda/shops/post_shops/main.main: " + err.Error()
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
