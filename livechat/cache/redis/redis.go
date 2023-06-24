// Package redis implements CacheClient for manipulating cache database.
package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// A Client contains redis database client
type Client struct {
	client *redis.Client // redis client representing a pool of zero or more underlying connections
}

// NewClient returns a new Client which contains redis client inside
func NewClient(addr string, password string) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	return &Client{client: rdb}
}

// Close close the redis client connection, releasing resources. Close return an error if it occurs.
func (c *Client) Close() error {
	err := c.client.Close()
	if err != nil {
		return err
	}
	return nil
}

// Set set the key and value with an expiration time. Set return an error if it occurs.
//
// Zero duration means that the key has no expiration time.
func (c *Client) Set(ctx context.Context, key string, value string, duration time.Duration) error {
	err := c.client.Set(ctx, key, value, time.Duration(duration)).Err()
	if err != nil {
		return err
	}
	return nil
}

// Get get the key and value with an expiration time. Get return an error if it occurs.
//
// Zero duration means that the key has no expiration time.
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

// GetShopConnections return all currently open connectionIDs of specific shop. GetShopConnections return an error if it occurs.
func (c *Client) GetShopConnections(ctx context.Context, shopID string) ([]string, error) {
	keys, err := c.client.Keys(ctx, shopID+":*").Result()
	if err != nil {
		return nil, err
	}
	connectionIDs := make([]string, len(keys))
	for i, key := range keys {
		connectionIDs[i] = key[len(shopID)+1:]
	}
	return connectionIDs, nil
}

// SetShopConnection set key value in the cache for a duration of time. SetShopConnection return an error if it occurs.
//
// set the key value as "<shopID>:<connectionID>":"<shopID>".
func (c *Client) SetShopConnection(ctx context.Context, shopID string, connectionID string, duration time.Duration) error {
	err := c.client.Set(ctx, shopID+":"+connectionID, shopID, time.Duration(duration)).Err()
	if err != nil {
		return err
	}
	return nil
}

// DeleteConnectionID delete an existing connectionID in the cache. DeleteConnectionID return an error if it occurs.
func (c *Client) DeleteConnectionID(ctx context.Context, connectionID string) error {
	allKeys, err := c.client.Keys(ctx, "*").Result()
	for _, key := range allKeys {
		if key[len(key)-len(connectionID):] == connectionID {
			err = c.client.Del(ctx, key).Err()
			if err != nil {
				return err
			}
			break
		}
	}

	if err != nil {
		return err
	}
	return nil
}
