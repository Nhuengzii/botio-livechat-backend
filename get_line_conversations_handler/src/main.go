package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (_ events.APIGatewayProxyResponse, err error) {
	defer func() {
		if err != nil {
			log.Println("get line conversations handler: " + err.Error())
			logToDiscord("get line conversations handler: " + err.Error())
		}
	}()
	pathParams := req.PathParameters
	pageID := pathParams["page_id"]
	dbc, err := newDBClient(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, err
	}
	defer dbc.Close(ctx)
	conversations, err := dbc.getConversationsOfPage(ctx, pageID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, err
	}
	if len(conversations) == 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Not Found",
		}, nil
	}
	returnConversations := returnConversations{conversations}
	returnBody, err := json.Marshal(returnConversations)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(returnBody),
	}, nil
}
