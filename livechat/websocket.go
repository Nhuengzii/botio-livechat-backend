package livechat

import "context"

type WebsocketClient interface {
	Send(ctx context.Context, connectionID string, message string) error
}
