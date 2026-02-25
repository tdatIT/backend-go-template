package orm

import (
	"database/sql"
	"log/slog"

	"github.com/tdatIT/backend-go/internal/domain/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// _gorm orm struct
type ormImpl struct {
	db    *gorm.DB
	sqlDB *sql.DB
}

func newGormInstance(c Config) (ORM, error) {
	dial := postgres.Open(c.DSN)

	gConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   c.TablePrefix,
			SingularTable: true,
		},
	}

	db, err := gorm.Open(dial, gConfig)
	if err != nil {
		return nil, err
	}

	if c.Debug {
		db = db.Debug()
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("error while getting sql db from gorm",
			slog.String("dsn", c.DSN),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	if c.MaxOpenConns != 0 {
		sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	}
	if c.ConnMaxLifetime != 0 {
		sqlDB.SetConnMaxLifetime(c.ConnMaxLifetime)
	}
	if c.MaxIdleConns != 0 {
		sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	}
	if c.ConnMaxIdleTime != 0 {
		sqlDB.SetConnMaxIdleTime(c.ConnMaxIdleTime)
	}

	// migration tables
	err = db.AutoMigrate(
		&models.User{},
		&models.Session{},
		&models.TaskGroup{},
		&models.Task{},
	)
	if err != nil {
		slog.Error("auto migrate failed", slog.Any("err", err))
		return nil, err
	}

	return &ormImpl{
		db:    db,
		sqlDB: sqlDB,
	}, nil
}

func (g *ormImpl) SqlDB() *sql.DB {
	return g.sqlDB
}

func (g *ormImpl) GormDB() *gorm.DB {
	return g.db
}

func (g *ormImpl) Close() error {
	return g.sqlDB.Close()
}
