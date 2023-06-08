package main

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/internal/fbutil/webhook"
	"github.com/Nhuengzii/botio-livechat-backend/internal/sqswrapper"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Lambda struct {
	config
}

func (l Lambda) handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Facebook websocket verify lambda handler")

	if request.HTTPMethod == "GET" {
		log.Println("GET method called")
		err := webhook.VerifyConnection(request.QueryStringParameters, l.FacebookWebhookVerificationString)
		if err != nil {
			log.Println(err)
			return events.APIGatewayProxyResponse{
				StatusCode: 401,
				Body:       "Unauthorized",
			}, err
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       request.QueryStringParameters["hub.challenge"],
		}, nil
	} else if request.HTTPMethod == "POST" {
		log.Println("POST method  called")
		start := time.Now()
		// new session

		// verify Signature
		err := webhook.VerifyMessageSignature(request.Headers, []byte(request.Body), l.FacebookAppSecret)
		if err != nil {
			log.Println(err)
			return events.APIGatewayProxyResponse{
				StatusCode: 401,
				Body:       "Unauthorized",
			}, err
		}

		msg := request.Body

		err = l.SqsClient.SendMessage(l.SqsQueueURL, msg)
		if err != nil {
			log.Println(err)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       "Internal Server Error",
			}, err
		}
		elasped := time.Since(start)
		log.Println("Elapsed : ", elasped)

		response = events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "OK",
		}

	} else {
		log.Printf("%v : method does not exist", request.HTTPMethod)
		return events.APIGatewayProxyResponse{}, errors.New("method does not exist")
	}
	return response, nil
}

func main() {
	l := Lambda{
		config: config{
			DiscordWebhookURL:                 os.Getenv("DISCORD_WEBHOOK_URL"),
			SqsQueueURL:                       os.Getenv("SQS_QUEUE_URL"),
			FacebookAppSecret:                 os.Getenv("FACEBOOK_APP_SECRET"), // TODO to be removed and get from some db instead
			FacebookWebhookVerificationString: os.Getenv("FACEBOOK_WEBHOOK_VERIFICATION_STRING"),
			SqsClient:                         *sqswrapper.NewClient(),
		},
	}
	lambda.Start(l.handler)
}
