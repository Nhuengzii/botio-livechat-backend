package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)

func main() {
	fmt.Println("Got connect")
	lambda.Start(Handler)
}

func Handler(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionID := request.RequestContext.ConnectionID
	log.Println("Livechat websocket connect lambda handler")
	log.Println("ConnectionID: ", connectionID)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}
