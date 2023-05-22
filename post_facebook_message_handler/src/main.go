package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(context context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	start := time.Now()
	discordLog(fmt.Sprint("-------Post-FacebookMessage-handler!!!!"))

	discordLog(fmt.Sprintf("Elasped : %v", time.Since(start)))
	return events.APIGatewayProxyResponse{}, nil
}
