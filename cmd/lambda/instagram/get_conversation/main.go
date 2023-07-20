package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/getconversation"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/apigateway"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	errNoShopIDPath         = errors.New("err path parameter parameters shop_id not given")
	errNoPageIDPath         = errors.New("err path parameter parameters page_id not given")
	errNoConversationIDPath = errors.New("err path parameter parameters conversation_id not given")
	errConversationNotExist = errors.New("err conversation does not exist")
)

func (c *config) handler(ctx context.Context, request events.APIGatewayProxyRequest) (_ events.APIGatewayProxyResponse, err error) {
	defer func() {
		if err != nil {
			discord.Log(c.discordWebhookURL, fmt.Sprintln(err))
		}
	}()

	//**path params checking//
	pathParams := request.PathParameters
	shopID, ok := pathParams["shop_id"]
	if !ok {
		return apigateway.NewProxyResponse(400, errNoShopIDPath.Error(), "*"), nil
	}
	pageID, ok := pathParams["page_id"]
	if !ok {
		return apigateway.NewProxyResponse(400, errNoPageIDPath.Error(), "*"), nil
	}
	conversationID, ok := pathParams["conversation_id"]
	if !ok {
		return apigateway.NewProxyResponse(400, errNoConversationIDPath.Error(), "*"), nil
	}
	//**end path params checking//

	stdConversation, err := c.dbClient.QueryConversation(ctx, shopID, pageID, conversationID)
	if err != nil {
		if errors.Is(err, mongodb.ErrNoDocuments) {
			return apigateway.NewProxyResponse(404, errConversationNotExist.Error(), "*"), nil
		}
		return apigateway.NewProxyResponse(502, "Bad Gateway", "*"), err
	}
	getConversationResponse := getconversation.Response{
		Conversation: stdConversation,
	}

	jsonBodyByte, err := json.Marshal(getConversationResponse)
	if err != nil {
		return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
	}
	return apigateway.NewProxyResponse(200, string(jsonBodyByte), "*"), nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2500*time.Millisecond)
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
