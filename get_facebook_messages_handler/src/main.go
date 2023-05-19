package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(context context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	discordLog("facebook get messages handler!!!")

	pathParams := request.PathParameters
	// shopID := pathParams["shop_id"]
	pageID := pathParams["page_id"]
	conversationID := pathParams["conversation_id"]

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}
