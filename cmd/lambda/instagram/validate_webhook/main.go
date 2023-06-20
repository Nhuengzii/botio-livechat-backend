package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/sqswrapper"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var errMethodNotAllowed = errors.New("HTTP method not allowed")

func (c *config) handler(ctx context.Context, request events.APIGatewayProxyRequest) (_ events.APIGatewayProxyResponse, err error) {
	defer func() {
		if err != nil {
			discord.Log(c.discordWebhookURL, fmt.Sprint(err))
		}
	}()

	if request.HTTPMethod == "GET" {
		err := VerifyConnection(request.QueryStringParameters, c.instagramWebhookVerificationString)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 401,
				Body:       "Unauthorized",
			}, err
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       request.QueryStringParameters["hub.challenge"],
		}, err
	} else if request.HTTPMethod == "POST" {
		// new session

		// verify Signature
		err := VerifyMessageSignature(request.Headers, []byte(request.Body), c.instagramAppSecret)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 401,
				Body:       "Unauthorized",
			}, err
		}

		msg := request.Body

		err = c.sqsClient.SendMessage(c.sqsQueueURL, msg)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       "Internal Server Error",
			}, err
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "OK",
		}, err

	} else {
		return events.APIGatewayProxyResponse{
			StatusCode: 405,
			Body:       "Method Not Allowed",
		}, errMethodNotAllowed
	}
}

func main() {
	var (
		discordWebhookURL                  = os.Getenv("DISCORD_WEBHOOK_URL")
		sqsQueueURL                        = os.Getenv("SQS_QUEUE_URL")
		appSecret                          = os.Getenv("APP_SECRET")
		instagramWebhookVerificationString = os.Getenv("INSTAGRAM_WEBHOOK_VERIFICATION_STRING")
		awsRegion                          = os.Getenv("AWS_REGION")
	)
	c := config{
		discordWebhookURL:                  discordWebhookURL,
		sqsQueueURL:                        sqsQueueURL,
		instagramAppSecret:                 appSecret, // TODO to be removed and get from some db instead
		instagramWebhookVerificationString: instagramWebhookVerificationString,
		sqsClient:                          sqswrapper.NewClient(awsRegion),
	}
	lambda.Start(c.handler)
}
