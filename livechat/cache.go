package livechat

import (
	"context"
	"time"
)

// CacheClient is an interface for cache client, for example, Redis
type CacheClient interface {
	// Close close the cache's client connection, releasing resources. Close return an error if it occurs.
	Close() error
	// Set set the key and value with an expiration time. Set return an error if it occurs.
	//
	// Zero duration means that the key has no expiration time.
	Set(ctx context.Context, key string, value string, duration time.Duration) error
	// Get get the key and value with an expiration time. Get return an error if it occurs.
	//
	// Zero duration means that the key has no expiration time.
	Get(ctx context.Context, key string) (string, error)
	// GetShopConnections return all currently open connectionIDs of specific shop. GetShopConnections return an error if it occurs.
	GetShopConnections(ctx context.Context, shopID string) ([]string, error)
	// SetShopConnection set key value in the cache for a duration of time. SetShopConnection return an error if it occurs.
	//
	// set the key value as "<shopID>:<connectionID>":"<shopID>".
	SetShopConnection(ctx context.Context, shopID string, connectionID string, duration time.Duration) error
	// DeleteConnectionID delete an existing connectionID in the cache. DeleteConnectionID return an error if it occurs.
	DeleteConnectionID(ctx context.Context, connectionID string) error
}
