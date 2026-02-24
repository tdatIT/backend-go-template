package cache

import (
	"context"
	"log/slog"
	"time"

	"github.com/tdatIT/backend-go/pkgs/db/rdclient"
)

type redisInstance struct {
	rdc rdclient.RedisClient
}

func NewRedisCacheClient(redisClient rdclient.RedisClient) Cache {
	cacheEngine := redisInstance{
		rdc: redisClient,
	}
	return &cacheEngine
}

func (r redisInstance) Get(ctx context.Context, key string) ([]byte, error) {
	result := r.rdc.Client().Get(ctx, key)
	val, err := result.Bytes()
	if err != nil {
		slog.Error("error while getting key from redis", slog.String("key", key), slog.Any("err", err))
		return nil, err
	}

	return val, err
}

func (r redisInstance) Set(ctx context.Context, key string, val interface{}, ttl time.Duration) error {
	result := r.rdc.Client().Set(ctx, key, val, ttl)
	return result.Err()
}

func (r redisInstance) Delete(ctx context.Context, key string) error {
	result := r.rdc.Client().Del(ctx, key)
	return result.Err()
}

func (r redisInstance) Expire(ctx context.Context, key string, ttl time.Duration) error {
	result := r.rdc.Client().Expire(ctx, key, ttl)
	return result.Err()
}
