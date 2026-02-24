package rdclient

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tdatIT/backend-go/config"
)

type RedisClient interface {
	Client() redis.UniversalClient
	Close() error
}

func NewRedisClient(cfg *config.ServiceConfig) RedisClient {
	if len(cfg.Redis.Address) == 0 {
		slog.Error("redis address list is empty")
		os.Exit(1)
	}

	connOpts := redis.UniversalOptions{
		Addrs:           cfg.Redis.Address,
		Username:        cfg.Redis.Username,
		Password:        cfg.Redis.Password,
		PoolSize:        cfg.Redis.PoolSize,
		DB:              cfg.Redis.DB,
		ConnMaxLifetime: cfg.Redis.ConnMaxLifeTime,
		ConnMaxIdleTime: cfg.Redis.ConnMaxIdleTime,
		MaxIdleConns:    cfg.Redis.MaxIdleConn,
		MinIdleConns:    cfg.Redis.MinIdleConn,
		DialTimeout:     5 * time.Second,
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    5 * time.Second,
		MaxRetries:      3,
	}

	if cfg.Redis.TLS.Enabled {
		tlsConfig, err := configureTLS(cfg)
		if err != nil {
			slog.Error("failed to configure TLS for Redis", slog.Any("err", err))
			os.Exit(1)
		}
		connOpts.TLSConfig = tlsConfig
	}

	switch cfg.Redis.Mode {
	case "cluster":
		// Cluster mode relies on multiple node addresses
		if len(connOpts.Addrs) == 0 {
			slog.Error("redis cluster mode requires at least one address")
			os.Exit(1)
		}
	case "sentinel":
		if cfg.Redis.MasterName == "" {
			slog.Error("redis sentinel mode requires master name")
			os.Exit(1)
		}
		connOpts.MasterName = cfg.Redis.MasterName
	default:
		// Treat any other value as standalone
		if len(connOpts.Addrs) == 0 {
			slog.Error("redis standalone mode requires a target address")
			os.Exit(1)
		}
		// Ensure only the primary address is used for standalone setups
		connOpts.Addrs = []string{connOpts.Addrs[0]}
	}

	client := redis.NewUniversalClient(&connOpts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		slog.Error("redis client ping failed", slog.Any("err", err))
		os.Exit(1)
	}

	slog.Info("redis client connected successfully")
	return &redisClient{
		_client: client,
	}
}

type redisClient struct {
	_client redis.UniversalClient
}

func (r *redisClient) Client() redis.UniversalClient {
	return r._client
}

func (r *redisClient) Close() error {
	return r._client.Close()
}
