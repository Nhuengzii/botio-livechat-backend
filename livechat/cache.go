package livechat

import "context"

type CacheClient interface {
	Close(ctx context.Context) error
	Set(ctx context.Context, key string, value string) error
	Get(ctx context.Context, key string) (string, error)
	GetShopConnections(ctx context.Context, shopID string) ([]string, error)
	DeleteConnectionID(ctx context.Context, shopID string, connectionID string) error
}
