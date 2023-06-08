package transport

import "github.com/aws/aws-lambda-go/events"

func SendError(statusCode int, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       body,
	}
}
