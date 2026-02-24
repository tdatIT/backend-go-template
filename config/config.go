package config

import (
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type ServiceConfig struct {
	Server   Server
	Database Database
	Redis    Redis
	Logger   Logger
	Auth     Auth
}

type Server struct {
	Name           string
	BuildVer       string
	HttpPort       string
	GrpcPort       string
	RequestTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	DefaultTimeout time.Duration
	DebugMode      bool
}

type Database struct {
	Host            string
	Port            int
	UserName        string
	Password        string
	Database        string
	Schema          string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	Options         string
}

type Redis struct {
	Mode            string
	Username        string
	Password        string
	PoolSize        int
	DB              int
	MasterName      string // for sentinel mode
	Address         []string
	MaxIdleConn     int
	MinIdleConn     int
	ConnMaxIdleTime time.Duration
	ConnMaxLifeTime time.Duration
	TLS             struct {
		Enabled            bool
		InsecureSkipVerify bool
		CertFilePath       string
	}
}

type Logger struct {
	Level string
}

type Auth struct {
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	GoogleClientID  string
}

type TelegramBot struct {
	Token  string
	ChatID string
}

// Get a config path for local or docker
func getDefaultConfig() string {
	return "/config/config"
}

// NewConfig Load config file from given path
func NewConfig() (*ServiceConfig, error) {
	config := ServiceConfig{}
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = getDefaultConfig()
	}
	log.Println("Loading config from path:", path)

	v := viper.New()
	v.SetConfigName(path)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		slog.Error("error reading config file", slog.String("err", err.Error()))
		return nil, err
	}

	err := v.Unmarshal(&config)
	if err != nil {
		slog.Error("error unmarshalling config file", slog.String("err", err.Error()))
		return nil, err
	}

	return &config, nil
}
