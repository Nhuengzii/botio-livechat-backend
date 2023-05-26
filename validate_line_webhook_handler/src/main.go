package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("webhook handler: new webhook received")
	logToDiscord("webhook handler: new webhook received")
	webhookBody := req.Body
	lineSignature := req.Headers["x-line-signature"]
	if err := validateSignature(lineChannelSecret, lineSignature, webhookBody); err != nil {
		log.Println("webhook handler: " + err.Error())
		logToDiscord("webhook handler: " + err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: 401,
			Body:       "Unauthorized",
		}, err
	}
	if err := sendSQSMessage(webhookBody); err != nil {
		log.Println("webhook handler: " + err.Error())
		logToDiscord("webhook handler: " + err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, err
	}
	log.Println("webhook handler: webhook body sent to sqs")
	logToDiscord("webhook handler: webhook body sent to sqs")
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "OK",
	}, nil
}
