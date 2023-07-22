package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/apigateway"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/getmessages"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	errSkipIntOnly  = errors.New("err skip parameter can only be integer")
	errLimitIntOnly = errors.New("err limit parameter can only be integer")
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

	messages := []stdmessage.StdMessage{}

	skipString, ok := req.QueryStringParameters["skip"]
	var skipPtr *int
	if skipString != "" {
		skip, err := strconv.Atoi(skipString)
		if err != nil {
			return apigateway.NewProxyResponse(400, errSkipIntOnly.Error(), "*"), nil
		}
		skipPtr = &skip
	}

	limitString, ok := req.QueryStringParameters["limit"]
	var limitPtr *int
	if limitString != "" {
		limit, err := strconv.Atoi(limitString)
		if err != nil {
			return apigateway.NewProxyResponse(400, errLimitIntOnly.Error(), "*"), nil
		}
		limitPtr = &limit
	}

	queryStringParameters := req.QueryStringParameters
	filterString, ok := queryStringParameters["filter"]
	if !ok {
		messages, err = c.dbClient.ListMessages(ctx, shopID, stdmessage.Platform("line"), pageID, conversationID, skipPtr, limitPtr)
	} else {
		filter := getmessages.Filter{}
		err = json.Unmarshal([]byte(filterString), &filter)
		if err != nil {
			return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
		}
		messages, err = c.dbClient.ListMessagesWithMessage(ctx, shopID, stdmessage.PlatformLine, pageID, conversationID, filter.Message, skipPtr, limitPtr)
	}
	if err != nil {
		return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
	}

	resp := getmessages.Response{
		Messages: messages,
	}
	responseJSON, err := json.Marshal(resp)
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
		CollectionTemplates:     "templates",
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
