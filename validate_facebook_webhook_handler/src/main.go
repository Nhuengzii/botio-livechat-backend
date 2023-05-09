package main

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// TODO: validate facebook post request
func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Facebook websocket verify lambda handler")
	var response events.APIGatewayProxyResponse

	if request.HTTPMethod == "GET" {
		log.Println("GET method called")
		err := VerificationCheck(request.QueryStringParameters)
		if err != nil {
			log.Println(err)
			return events.APIGatewayProxyResponse{}, err
		}
		response = events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       request.QueryStringParameters["hub.challenge"],
		}
	} else if request.HTTPMethod == "POST" {
		log.Println("POST method called")
		start := time.Now()
		//new session
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String("ap-southeast-1"),
		})
		if err != nil {
			log.Println(err)
			return events.APIGatewayProxyResponse{}, err
		}

		sqsClient := sqs.New(sess)
		queueUrl := os.Getenv("SQS_QUEUE_URL")
		msg := request.Body

		log.Println(msg)
		SendQueueMessage(msg, sqsClient, queueUrl)
		if err != nil {
			log.Println(err)
			return events.APIGatewayProxyResponse{}, err
		}
		elasped := time.Since(start)
		log.Println("Elapsed : ", elasped)

		response = events.APIGatewayProxyResponse{
			StatusCode: 200,
		}

	} else {
		log.Printf("%v : method does not exist", request.HTTPMethod)
		return events.APIGatewayProxyResponse{}, errors.New("method does not exist")
	}

	return response, nil

}
