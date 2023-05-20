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

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	start := time.Now()
	discordLog("get_facebook_conversation handler!!!")

	pathParams := request.PathParameters
	// shopID := pathParams["shop_id"]
	pageID := pathParams["page_id"]

	var outputMessage OutputMessage
	err := QueryConversations(pageID, &outputMessage)
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
