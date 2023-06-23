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

var errTwoFilterParamsInOneRequest = errors.New("err path parameters filter can only give 1 filter per 1 request")

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

	skipString := req.QueryStringParameters["skip"]
	var skipPtr *int
	if skipString != "" {
		skip, err := strconv.Atoi(skipString)
		if err != nil {
			return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
		}
		skipPtr = &skip
	}

	limitString := req.QueryStringParameters["limit"]
	var limitPtr *int
	if limitString != "" {
		limit, err := strconv.Atoi(limitString)
		if err != nil {
			return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
		}
		limitPtr = &limit
	}

	conversations := []stdconversation.StdConversation{}

	filterQueryString, ok := req.QueryStringParameters["filter"]
	if !ok { // no need to query with filter
		conversations, err = c.dbClient.ListConversationsOfAllPlatformsOfShop(ctx, shopID, skipPtr, limitPtr)
		if err != nil {
			return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
		}
	} else { // need to query with filter
		var filter getconversations.Filter

		err := json.Unmarshal([]byte(filterQueryString), &filter)
		if err != nil {
			return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
		}

		if filter.Message != "" && filter.ParticipantsUsername != "" {
			return apigateway.NewProxyResponse(400, "Bad Request", "*"), errTwoFilterParamsInOneRequest
		} else if filter.ParticipantsUsername != "" { // query with ParticipantsUsername
			conversations, err = c.dbClient.QueryConversationsOfAllPlatformWithParticipantsName(ctx, shopID, filter.ParticipantsUsername, skipPtr, limitPtr)
			if err != nil {
				return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
			}
		} else if filter.Message != "" { // query with message
			conversations, err = c.dbClient.QueryConversationsOfAllPlatformWithMessage(ctx, shopID, filter.Message, skipPtr, limitPtr)
			if err != nil {
				return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
			}
		}
	}

	response := getconversations.Response{
		Conversations: conversations,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		discord.Log(c.discordWebhookURL, "json Marshal error")
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
