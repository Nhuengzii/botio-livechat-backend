package apigateway

import "github.com/aws/aws-lambda-go/events"

func NewProxyResponse(statusCode int, body string, allowOrigin string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       body,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": allowOrigin,
		},
	}
}
