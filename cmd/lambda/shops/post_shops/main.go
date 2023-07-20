package main

import (
	"context"
	"encoding/json"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/postshop"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/shopcfg"
	"github.com/google/uuid"
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
		return apigateway.NewProxyResponse(400, "Bad Request: Request body must not be empty.", "*"), nil
	}

	postShopReq := postshop.Request{}
	err = json.Unmarshal([]byte(reqBody), &postShopReq)
	if err != nil {
		return apigateway.NewProxyResponse(400, "Bad Request: Check request body.", "*"), nil
	}

	newShopID := uuid.New().String()
	newShop := shops.Shop{
		ShopID:         newShopID,
		FacebookPageID: postShopReq.FacebookPageID,
		FacebookAuthentication: shops.FacebookAuthentication{
			AccessToken: postShopReq.FacebookAccessToken,
		},
		InstagramPageID: postShopReq.InstagramPageID,
		InstagramAuthentication: shops.InstagramAuthentication{
			AccessToken: postShopReq.InstagramAccessToken,
		},
		LinePageID: postShopReq.LinePageID,
		LineAuthentication: shops.LineAuthentication{
			AccessToken: postShopReq.LineAccessToken,
			Secret:      postShopReq.LineSecret,
		},
	}
	err = c.dbClient.InsertShop(ctx, newShop)
	if err != nil {
		return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
	}

	newShopConfig := shopcfg.Config{
		ShopID: newShopID,
	}
	err = c.dbClient.InsertShopConfig(ctx, newShopConfig)
	if err != nil {
		return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
	}

	resp := postshop.Response{
		ShopID: newShopID,
	}
	respJSON, err := json.Marshal(resp)
	if err != nil {
		return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
	}
	return apigateway.NewProxyResponse(200, string(respJSON), "*"), nil
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
		CollectionShopConfig:    "shop_config",
		CollectionTemplates:     "templates",
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
