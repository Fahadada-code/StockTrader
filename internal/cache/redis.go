package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	Client *redis.Client
}

func NewRedisCache(addr string) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &RedisCache{Client: rdb}
}

func (rc *RedisCache) SetSnapshot(ctx context.Context, symbol string, data interface{}, expiration time.Duration) error {
	return rc.Client.Set(ctx, "snapshot:"+symbol, data, expiration).Err()
}

func (rc *RedisCache) GetSnapshot(ctx context.Context, symbol string) (string, error) {
	return rc.Client.Get(ctx, "snapshot:"+symbol).Result()
}

func (rc *RedisCache) IncrementSubscriberCount(ctx context.Context, symbol string) error {
	return rc.Client.Incr(ctx, "subs:"+symbol).Err()
}

func (rc *RedisCache) DecrementSubscriberCount(ctx context.Context, symbol string) error {
	return rc.Client.Decr(ctx, "subs:"+symbol).Err()
}

func (rc *RedisCache) GetHotSymbols(ctx context.Context, limit int64) ([]string, error) {
	// Simple implementation using keys, better would be a sorted set
	// For now, just return anything with subs > 0
	return rc.Client.Keys(ctx, "subs:*").Result()
}
