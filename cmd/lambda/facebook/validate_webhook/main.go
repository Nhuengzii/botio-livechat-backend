package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

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

	discord.Log(c.discordWebhookURL, "Facebook webhook verify lambda handler")

	if request.HTTPMethod == "GET" {
		err := VerifyConnection(request.QueryStringParameters, c.facebookWebhookVerificationString)
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
		start := time.Now()
		// new session

		// verify Signature
		err := VerifyMessageSignature(request.Headers, []byte(request.Body), c.facebookAppSecret)
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
				StatusCode: 502,
				Body:       "Bad Gateway",
			}, err
		}
		elapsed := time.Since(start)
		discord.Log(c.discordWebhookURL, fmt.Sprint("elapsed : ", elapsed))
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "OK",
		}, err

	} else {
		return events.APIGatewayProxyResponse{
			StatusCode: 405,
			Body:       "Method",
		}, errMethodNotAllowed
	}
}

func main() {
	var (
		discordWebhookURL                 = os.Getenv("DISCORD_WEBHOOK_URL")
		sqsQueueURL                       = os.Getenv("SQS_QUEUE_URL")
		appSecret                         = os.Getenv("APP_SECRET")
		facebookWebhookVerificationString = os.Getenv("FACEBOOK_WEBHOOK_VERIFICATION_STRING")
		awsRegion                         = os.Getenv("AWS_REGION")
	)
	c := config{
		discordWebhookURL:                 discordWebhookURL,
		sqsQueueURL:                       sqsQueueURL,
		facebookAppSecret:                 appSecret, // TODO to be removed and get from some db instead
		facebookWebhookVerificationString: facebookWebhookVerificationString,
		sqsClient:                         sqswrapper.NewClient(awsRegion),
	}
	lambda.Start(c.handler)
}
