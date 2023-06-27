// Package websocketwrapper implements apigateway's websocket for manipulating websocket operations.
package websocketwrapper

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
)

// A Client contains apigatewaymanagementapi's client.
type Client struct {
	client *apigatewaymanagementapi.Client // used to perform various websocket operations
}

// NewClient create and return an apigateway's websocket client.
func NewClient(ctx context.Context, endpoint string) *Client {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-southeast-1"))
	if err != nil {
		return nil
	}
	svc := apigatewaymanagementapi.NewFromConfig(cfg, func(o *apigatewaymanagementapi.Options) {
		o.EndpointResolver = apigatewaymanagementapi.EndpointResolverFunc(func(region string, options apigatewaymanagementapi.EndpointResolverOptions) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           endpoint,
				SigningRegion: region,
			}, nil
		})
	})
	return &Client{client: svc}
}

// Send recieve a message string and send the message into a specific websocket connectionID.
// Return an error if it occurs.
func (c *Client) Send(ctx context.Context, connectionID string, message string) error {
	input := &apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(connectionID), Data: []byte(message)}
	_, err := c.client.PostToConnection(context.Background(), input)
	if err != nil {
		return err
	}
	return nil
}
