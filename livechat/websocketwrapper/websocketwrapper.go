package websocketwrapper

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
)

type Client struct {
	client *apigatewaymanagementapi.Client
}

func NewClient(enpoint string) *Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))
	if err != nil {
		panic(err)
	}
	client := apigatewaymanagementapi.NewFromConfig(cfg, func(o *apigatewaymanagementapi.Options) {
		o.EndpointResolver = apigatewaymanagementapi.EndpointResolverFunc(func(region string, options apigatewaymanagementapi.EndpointResolverOptions) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           enpoint,
				SigningRegion: region,
			}, nil
		})
	})
	return &Client{client}
}

func (c *Client) Send(ctx context.Context, connectionID string, message string) error {
	input := &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(connectionID),
		Data:         []byte(message),
	}
	_, err := c.client.PostToConnection(ctx, input)
	if err != nil {
		return err
	}
	return nil
}
