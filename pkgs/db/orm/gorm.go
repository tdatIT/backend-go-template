package orm

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// ORM defines a interface for access the db.
type ORM interface {
	GormDB() *gorm.DB
	SqlDB() *sql.DB
	Close() error
}

// Config GORM Config
type Config struct {
	Debug           bool
	DBType          string
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	TablePrefix     string
}
