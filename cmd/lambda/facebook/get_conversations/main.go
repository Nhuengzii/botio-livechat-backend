package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/getconversations"
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

	discord.Log(c.discordWebhookURL, "facebook get conversations handler")

	pathParams := request.PathParameters
	shopID, ok := pathParams["shop_id"]
	if !ok {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Bad Request",
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, errNoShopIDPath
	}
	pageID, ok := pathParams["page_id"]
	if !ok {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Bad Request",
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, errNoPageIDPath
	}

	filterQueryString, ok := request.QueryStringParameters["filter"]
	stdConversations := []stdconversation.StdConversation{}
	if !ok { // no need to query
		discord.Log(c.discordWebhookURL, "no need to query")
		stdConversations, err = c.dbClient.QueryConversations(ctx, shopID, pageID)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       "Internal Server Error",
				Headers: map[string]string{
					"Access-Control-Allow-Origin": "*",
				},
			}, err
		}
	} else { // need to query
		var filter getconversations.Filter
		discord.Log(c.discordWebhookURL, filterQueryString)
		err := json.Unmarshal([]byte(filterQueryString), &filter)
		discord.Log(c.discordWebhookURL, fmt.Sprintf("%+v", filter))
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       "Internal Server Error",
				Headers: map[string]string{
					"Access-Control-Allow-Origin": "*",
				},
			}, err
		}
		if filter.Message != "" && filter.ParticipantsUsername != "" {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Bad Request",
				Headers: map[string]string{
					"Access-Control-Allow-Origin": "*",
				},
			}, errTwoFilterParamsInOneRequest
		} else if filter.ParticipantsUsername != "" { // query with ParticipantsUsername
			stdConversations, err = c.dbClient.QueryConversationsWithParticipantsName(ctx, shopID, stdconversation.PlatformFacebook, pageID, filter.ParticipantsUsername)
			if err != nil {
				return events.APIGatewayProxyResponse{
					StatusCode: 500,
					Body:       "Internal Server Error",
					Headers: map[string]string{
						"Access-Control-Allow-Origin": "*",
					},
				}, err
			}
		}
	}
	getConversationsResponse := getconversations.Response{
		Conversations: stdConversations,
	}

	jsonBodyByte, err := json.Marshal(getConversationsResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonBodyByte),
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
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
