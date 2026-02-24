package htlcheck

import (
	"context"
	"log/slog"
	"time"

	"github.com/hellofresh/health-go/v5"
	"github.com/tdatIT/backend-go/config"
	"github.com/tdatIT/backend-go/pkgs/db/orm"
	"github.com/tdatIT/backend-go/pkgs/db/rdclient"
)

// NewHealthCheckService creates a new health check service with database and redis checks
func NewHealthCheckService(
	cfg *config.ServiceConfig,
	_db orm.ORM,
	_redis rdclient.RedisClient,
) (*health.Health, error) {
	h, err := health.New(health.WithComponent(health.Component{
		Name:    cfg.Server.Name,
		Version: cfg.Server.BuildVer,
	}))
	if err != nil {
		return nil, err
	}

	// Database health check
	err = h.Register(health.Config{
		Name:      "postgres",
		Timeout:   time.Second * 5,
		SkipOnErr: false,
		Check: func(ctx context.Context) error {
			if err := _db.SqlDB().PingContext(ctx); err != nil {
				slog.Error("database health check failed", slog.Any("error", err))
				return err
			}
			return nil
		},
	})
	if err != nil {
		return nil, err
	}

	// Redis health check
	err = h.Register(health.Config{
		Name:      "redis",
		Timeout:   time.Second * 5,
		SkipOnErr: false,
		Check: func(ctx context.Context) error {
			if err := _redis.Client().Ping(ctx).Err(); err != nil {
				slog.Error("redis health check failed", slog.Any("error", err))
				return err
			}
			return nil
		},
	})
	if err != nil {
		return nil, err
	}

	return h, nil
}
