package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
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
		}
	}
	for _, conversation := range outputMessage.Conversations {
		discordLog(fmt.Sprintf("Last message in conversation %v is %v", conversation.ConversationID, conversation.LastActivity))
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}
}
