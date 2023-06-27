package livechat

import "context"

// WebsocketClient is an interface for websocket client, for example, API Gateway
type WebsocketClient interface {
	// Send recieve a message string and send the message into a specific websocket connection.
	// Return an error if it occurs.
	Send(ctx context.Context, connectionID string, message string) error
}
