package livechat

import "context"

// WebsocketClient is an interface for websocket client, for example, API Gateway
type WebsocketClient interface {
	Send(ctx context.Context, connectionID string, message string) error
}
