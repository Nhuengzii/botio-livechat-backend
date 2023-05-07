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
	}

	return response, nil

}
