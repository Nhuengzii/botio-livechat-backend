package livechat

import (
	"context"
	"time"
)

type CacheClient interface {
	Close() error
	Set(ctx context.Context, key string, value string, duration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	GetShopConnections(ctx context.Context, shopID string) ([]string, error)
	SetShopConnection(ctx context.Context, shopID string, connectionID string, duration time.Duration) error
	DeleteConnectionID(ctx context.Context, connectionID string) error
}
