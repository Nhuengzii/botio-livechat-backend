package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/patchshopcfg"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/templates"
	"github.com/google/uuid"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/apigateway"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func (c *config) handler(ctx context.Context, req events.APIGatewayProxyRequest) (_ events.APIGatewayProxyResponse, err error) {
	defer func() {
		if err != nil {
			logMessage := "cmd/lambda/shops/post_templates/main.config.handler: " + err.Error()
			log.Println(logMessage)
			discord.Log(c.discordWebhookURL, logMessage)
		}
	}()

	pathParameters := req.PathParameters
	shopID, ok := pathParameters["shop_id"]
	if !ok {
		return apigateway.NewProxyResponse(400, "BadRequest: shop_id must not be empty.", "*"), nil
	}

	reqBody := req.Body
	if reqBody == "" {
		return apigateway.NewProxyResponse(400, "Bad Request: Request body must not be empty.", "*"), nil
	}

	patchShopCfgRequest := patchshopcfg.Request{}
	err = json.Unmarshal([]byte(reqBody), &patchShopCfgRequest)
	if err != nil {
		return apigateway.NewProxyResponse(400, "Bad Request: Check request body.", "*"), nil
	}

	newTemplateID := uuid.New().String()
	newTemplate := templates.Template{
		ShopID:  shopID,
		ID:      newTemplateID,
		Payload: patchShopCfgRequest.TemplatePayload,
	}
	err = c.dbClient.AddShopNewTemplateMessage(ctx, newTemplate)
	if err != nil {
		if errors.Is(err, mongodb.ErrNoDocuments) {
			return apigateway.NewProxyResponse(404, "Not Found: Shop not found.", "*"), nil
		}
		return apigateway.NewProxyResponse(500, "Internal Server Error: ", "*"), nil
	}

	response := patchshopcfg.Response{
		TemplateID: newTemplateID,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return apigateway.NewProxyResponse(500, "Internal Server Error: ", "*"), nil
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
		CollectionShopConfig:    "shop_config",
		CollectionTemplates:     "templates",
	})
	if err != nil {
		logMessage := "cmd/lambda/shops/post_templates/main.main: " + err.Error()
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
