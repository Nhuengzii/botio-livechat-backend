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

	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/getconversations"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/apigateway"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdconversation"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	errNoPageIDPath                = errors.New("err path parameter parameters page_id not given")
	errNoShopIDPath                = errors.New("err path parameter parameters shop_id not given")
	errTwoFilterParamsInOneRequest = errors.New("err path parameters filter can only give 1 filter per 1 request")
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
		return apigateway.NewProxyResponse(400, "Bad Request", "*"), errNoShopIDPath
	}
	pageID, ok := pathParams["page_id"]
	if !ok {
		return apigateway.NewProxyResponse(400, "Bad Request", "*"), errNoPageIDPath
	}

	stdConversations := []stdconversation.StdConversation{}

	skipString, ok := request.QueryStringParameters["skip"]
	var skipPtr *int
	if skipString != "" {
		skip, err := strconv.Atoi(skipString)
		if err != nil {
			return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
		}
		skipPtr = &skip
	}

	limitString, ok := request.QueryStringParameters["limit"]
	var limitPtr *int
	if limitString != "" {
		limit, err := strconv.Atoi(limitString)
		if err != nil {
			return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
		}
		limitPtr = &limit
	}

	filterQueryString, ok := request.QueryStringParameters["filter"]
	if !ok { // no need to query with filter
		stdConversations, err = c.dbClient.QueryConversations(ctx, shopID, pageID, skipPtr, limitPtr)
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
			stdConversations, err = c.dbClient.QueryConversationsWithParticipantsName(ctx, shopID, stdconversation.PlatformFacebook, pageID, filter.ParticipantsUsername, skipPtr, limitPtr)
			if err != nil {
				return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
			}
		} else if filter.Message != "" { // query with message
			stdConversations, err = c.dbClient.QueryConversationsWithMessage(ctx, shopID, stdconversation.PlatformFacebook, pageID, filter.Message, skipPtr, limitPtr)
			if err != nil {
				return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
			}
		}
	}
	getConversationsResponse := getconversations.Response{
		Conversations: stdConversations,
	}

	jsonBodyByte, err := json.Marshal(getConversationsResponse)
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
