package main

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/fbutil/webhook"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/sqswrapper"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Lambda struct {
	config
}

// aliasing the types to keep lines short
type (
	Request  = events.APIGatewayProxyRequest
	Response = events.APIGatewayProxyResponse
)

func (l Lambda) handler(ctx context.Context, request Request) (Response, error) {
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
		}, err
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
				StatusCode: 502,
				Body:       "Bad Gateway",
			}, err
		}
		elasped := time.Since(start)
		log.Println("Elapsed : ", elasped)

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "OK",
		}, err

	} else {
		log.Printf("%v : method does not exist", request.HTTPMethod)
		return events.APIGatewayProxyResponse{
			StatusCode: 405,
			Body:       "Method",
		}, errors.New("Method Not Allowed")
	}
}

func main() {
	l := Lambda{
		config: config{
			DiscordWebhookURL:                 os.Getenv("DISCORD_WEBHOOK_URL"),
			SqsQueueURL:                       os.Getenv("SQS_QUEUE_URL"),
			FacebookAppSecret:                 os.Getenv("FACEBOOK_APP_SECRET"), // TODO to be removed and get from some db instead
			FacebookWebhookVerificationString: os.Getenv("FACEBOOK_WEBHOOK_VERIFICATION_STRING"),
			SqsClient:                         sqswrapper.NewClient(os.Getenv("AWS_REGION")),
		},
	}
	lambda.Start(l.handler)
}
