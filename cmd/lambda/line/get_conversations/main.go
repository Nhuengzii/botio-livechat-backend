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
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdconversation"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/getconversations"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	errSkipIntOnly                 = errors.New("err skip parameter can only be integer")
	errLimitIntOnly                = errors.New("err limit parameter can only be integer")
	errTwoFilterParamsInOneRequest = errors.New("err path parameters filter can only give 1 filter per 1 request")
)

func (c *config) handler(ctx context.Context, req events.APIGatewayProxyRequest) (_ events.APIGatewayProxyResponse, err error) {
	defer func() {
		if err != nil {
			logMessage := "cmd/lambda/line/get_conversations/main.config.handler: " + err.Error()
			log.Println(logMessage)
			discord.Log(c.discordWebhookURL, logMessage)
		}
	}()
	pathParameters := req.PathParameters
	shopID := pathParameters["shop_id"]
	pageID := pathParameters["page_id"]

	conversations := []stdconversation.StdConversation{}

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
		conversations, err = c.dbClient.ListConversations(ctx, shopID, stdconversation.PlatformLine, pageID, skipPtr, limitPtr)
	} else {
		filter := getconversations.Filter{}
		err = json.Unmarshal([]byte(filterString), &filter)
		if err != nil {
			return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
		}
		if (filter.ParticipantsUsername != "") && (filter.Message == "") {
			conversations, err = c.dbClient.ListConversationsWithParticipantsName(ctx, shopID, stdconversation.PlatformLine, pageID, filter.ParticipantsUsername, skipPtr, limitPtr)
		} else if (filter.ParticipantsUsername == "") && (filter.Message != "") {
			conversations, err = c.dbClient.ListConversationsWithMessage(ctx, shopID, stdconversation.PlatformLine, pageID, filter.Message, skipPtr, limitPtr)
		} else {
			return apigateway.NewProxyResponse(400, errTwoFilterParamsInOneRequest.Error(), "*"), nil
		}
	}
	if err != nil {
		return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
	}

	resp := getconversations.Response{
		Conversations: conversations,
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
		logMessage := "cmd/lambda/line/get_conversations/main.main: " + err.Error()
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
