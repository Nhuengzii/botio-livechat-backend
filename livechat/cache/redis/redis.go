package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	client *redis.Client
}

func NewClient(addr string, password string) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	return &Client{client: rdb}
}

func (c *Client) Close() error {
	err := c.client.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Set(ctx context.Context, key string, value string, duration int64) error {
	err := c.client.Set(ctx, key, value, time.Duration(duration)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

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

func (c *Client) SetShopConnection(ctx context.Context, shopID string, connectionID string, duration int64) error {
	err := c.client.Set(ctx, shopID+":"+connectionID, shopID, time.Duration(duration)).Err()
	if err != nil {
		return err
	}
	return nil
}

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
