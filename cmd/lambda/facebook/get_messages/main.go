package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/getmessages"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/apigateway"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	errNoShopIDPath         = errors.New("err path parameter parameters shop_id not given")
	errNoPageIDPath         = errors.New("err path parameter parameters page_id not given")
	errNoConversationIDPath = errors.New("err path parameter conversation_id not given")
	errSkipIntOnly          = errors.New("err skip parameter can only be integer")
	errLimitIntOnly         = errors.New("err limit parameter can only be integer")
)

func (c *config) handler(ctx context.Context, request events.APIGatewayProxyRequest) (_ events.APIGatewayProxyResponse, err error) {
	defer func() {
		if err != nil {
			discord.Log(c.discordWebhookURL, fmt.Sprintln(err))
		}
	}()

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

	stdMessages := []stdmessage.StdMessage{}

	skipString, ok := request.QueryStringParameters["skip"]
	var skipPtr *int
	if skipString != "" {
		skip, err := strconv.Atoi(skipString)
		if err != nil {
			return apigateway.NewProxyResponse(400, errSkipIntOnly.Error(), "*"), nil
		}
		skipPtr = &skip
	}

	limitString, ok := request.QueryStringParameters["limit"]
	var limitPtr *int
	if limitString != "" {
		limit, err := strconv.Atoi(limitString)
		if err != nil {
			return apigateway.NewProxyResponse(400, errLimitIntOnly.Error(), "*"), nil
		}
		limitPtr = &limit
	}

	filterQueryString, ok := request.QueryStringParameters["filter"]
	if !ok {
		stdMessages, err = c.dbClient.ListMessages(ctx, shopID, stdmessage.PlatformFacebook, pageID, conversationID, skipPtr, limitPtr)
		if err != nil {
			return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
		}
	} else {
		var filter getmessages.Filter
		err := json.Unmarshal([]byte(filterQueryString), &filter)

		stdMessages, err = c.dbClient.ListMessagesWithMessage(ctx, shopID, stdmessage.PlatformFacebook, pageID, conversationID, filter.Message, skipPtr, limitPtr)
		if err != nil {
			return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
		}
	}
	getMessagesResponse := getmessages.Response{
		Messages: stdMessages,
	}

	jsonBodyByte, err := json.Marshal(getMessagesResponse)
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
