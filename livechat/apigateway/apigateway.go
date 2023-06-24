// Package apigateway contains apigateway helper functions
package apigateway

import "github.com/aws/aws-lambda-go/events"

// NewProxyResponse returns an apigateway proxy response.
// It configures the response to be returned by API Gateway for the request.
//
// Should the caller wants to send json message they should marshal and stringify the body first.
//
// *** This function should only be used by lambda handlers ***
func NewProxyResponse(statusCode int, body string, allowOrigin string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       body,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": allowOrigin,
		},
	}
}
