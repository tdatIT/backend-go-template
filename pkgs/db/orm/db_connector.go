package orm

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/tdatIT/backend-go/config"
)

func NewDBConnection(config *config.ServiceConfig) ORM {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s %s",
		config.Database.Host,
		config.Database.Port,
		config.Database.UserName,
		config.Database.Password,
		config.Database.Database,
		config.Database.Options,
	)

	if config.Database.Schema != "" {
		dsn += fmt.Sprintf(" search_path=%s", config.Database.Schema)
	}

	cfg := Config{
		DSN:             dsn,
		MaxOpenConns:    config.Database.MaxOpenConns,
		MaxIdleConns:    config.Database.MaxIdleConns,
		ConnMaxLifetime: config.Database.ConnMaxLifetime,
		ConnMaxIdleTime: config.Database.ConnMaxIdleTime,
		Debug:           config.Server.DebugMode,
	}
	conn, err := newGormInstance(cfg)
	if err != nil {
		slog.Error("error while creating db connection", slog.Any("err", err))
		os.Exit(1)
	}

	slog.Info("db connection established",
		slog.String("host", fmt.Sprintf("%v:%v", config.Database.Host, config.Database.Port)),
		slog.String("db", config.Database.Database),
		slog.String("schema", config.Database.Schema),
	)

	return conn
}
