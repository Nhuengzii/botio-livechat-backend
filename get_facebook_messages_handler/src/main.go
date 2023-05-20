package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(context context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	discordLog("facebook get messages handler!!!")
	start := time.Now()

	pathParams := request.PathParameters
	// shopID := pathParams["shop_id"]
	pageID := pathParams["page_id"]
	conversationID := pathParams["conversation_id"]

	var outputMessage OutputMessage
	err := QueryMessages(pageID, conversationID, &outputMessage)
	if err != nil {
		discordLog(fmt.Sprint(err))
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadGateway,
		}, err
	}

	jsonBodyByte, err := json.Marshal(outputMessage)
	jsonString := string(jsonBodyByte)
	if err != nil {
		discordLog(fmt.Sprint(err))
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadGateway,
		}, err
	}

	discordLog(fmt.Sprintf("Total Elasped time: %v", time.Since(start)))

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       jsonString,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}
