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
			log.Println("get line message handler: " + err.Error())
			logToDiscord("get line message handler: " + err.Error())
		}
	}()
	pathParams := req.PathParameters
	// conversationID is designed to be unique (conversationID=botUserID:userID)
	// so other path parameters are not needed
	conversationID := pathParams["conversation_id"]
	logToDiscord("get line message handler: conversationID=" + conversationID)
	dbc, err := newDBClient(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, err
	}
	defer dbc.Close(ctx)
	messages, err := dbc.getMessagesInConversation(ctx, conversationID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, err
	}
	if len(messages) == 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Not Found",
		}, nil
	}
	returnMessages := returnMessages{messages}
	returnBody, err := json.Marshal(returnMessages)
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
