package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) {
	discordLog("get_facebook_conversation handler!!!")

	pathParams := request.PathParameters
	// shopID := pathParams["shop_id"]
	pageID := pathParams["page_id"]

	var outputMessage OutputMessage
	err := QueryConversations(pageID, &outputMessage)
	if err != nil {
		discordLog(fmt.Sprint(err))
	}
}
