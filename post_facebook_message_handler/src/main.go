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

func handler(context context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	discordLog(fmt.Sprint("-------Post-FacebookMessage-handler!!!!"))
	return events.APIGatewayProxyResponse{}, nil
}
