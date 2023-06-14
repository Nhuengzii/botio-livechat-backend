package websocketwrapper

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/aws"
)

func (c *Client) Send(ctx context.Context, connectionID string, message string) error {
	input := &apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(connectionID), Data: []byte(message)}
	_, err := c.client.PostToConnection(context.Background(), input)
	if err != nil {
		return err
	}
	return nil
}
