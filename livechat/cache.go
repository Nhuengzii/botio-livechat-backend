package livechat

import "context"

type CacheClient interface {
	Close() error
	Set(ctx context.Context, key string, value string, duration int64) error
	Get(ctx context.Context, key string) (string, error)
	GetShopConnections(ctx context.Context, shopID string) ([]string, error)
	SetShopConnection(ctx context.Context, shopID string, connectionID string, duration int64) error
	DeleteConnectionID(ctx context.Context, connectionID string) error
}
