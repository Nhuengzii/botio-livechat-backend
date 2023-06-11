package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/response/getconversationsresp"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"os"
)

func (c *config) handler(ctx context.Context, req events.APIGatewayProxyRequest) (_ events.APIGatewayProxyResponse, err error) {
	defer func() {
		if err != nil {
			logMessage := "lambda/line/get_conversations/main.config.handler: " + err.Error()
			log.Println(logMessage)
			discord.Log(c.discordWebhookURL, logMessage)
		}
	}()
	if c.dbClient == nil {
		c.dbClient, err = mongodb.NewClient(ctx, &mongodb.Target{
			URI:                     c.mongodbURI,
			Database:                c.mongodbDatabase,
			CollectionConversations: c.mongodbCollectionLineConversations,
		})
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Headers: map[string]string{
					"Access-Control-Allow-Origin": "*",
				},
				Body: "Internal Server Error",
			}, err
		}
	}
	pathParameters := req.PathParameters
	pageID := pathParameters["page_id"]
	conversations, err := c.dbClient.QueryConversations(ctx, pageID)
	if err != nil {
		if errors.Is(err, mongodb.ErrNoConversations) {
			return events.APIGatewayProxyResponse{
				StatusCode: 404,
				Headers: map[string]string{
					"Access-Control-Allow-Origin": "*",
				},
				Body: "Not Found",
			}, nil
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
			Body: "Internal Server Error",
		}, nil
	}
	resp := &getconversationsresp.Resp{
		Conversations: conversations,
	}
	responseJSON, err := json.Marshal(resp)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
			Body: "Internal Server Error",
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
		Body: string(responseJSON),
	}, nil
}

func main() {
	c := &config{
		discordWebhookURL:                  os.Getenv("DISCORD_WEBHOOK_URL"),
		mongodbURI:                         os.Getenv("MONGODB_URI"),
		mongodbDatabase:                    os.Getenv("MONGODB_DATABASE"),
		mongodbCollectionLineConversations: os.Getenv("MONGODB_COLLECTION_LINE_CONVERSATIONS"),
		dbClient:                           nil,
	}
	lambda.Start(c.handler)
}
