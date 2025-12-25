package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"flexi-redirector/internal/env"
	"flexi-redirector/internal/features"
	"flexi-redirector/internal/features/countviews"
)

type ServerConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type DBConfig struct {
	Driver     string // postgres|sqlite
	DBURL      string // for postgres
	SQLitePath string // for sqlite file path or DSN

	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration

	AutoMigrate bool
}

type Config struct {
	Server ServerConfig
	DB     DBConfig

	CountViews *countviews.Feature
}

func LoadFromEnv() (Config, error) {
	var cfg Config

	cfg.Server.Port = env.String("SERVER_PORT", "8080")
	cfg.Server.ReadTimeout = env.Duration("SERVER_READ_TIMEOUT", 10*time.Second)
	cfg.Server.WriteTimeout = env.Duration("SERVER_WRITE_TIMEOUT", 10*time.Second)
	cfg.Server.IdleTimeout = env.Duration("SERVER_IDLE_TIMEOUT", 120*time.Second)
	cfg.Server.ShutdownTimeout = env.Duration("SERVER_SHUTDOWN_TIMEOUT", 5*time.Second)

	cfg.DB.Driver = strings.ToLower(env.String("DB_DRIVER", "postgres"))
	cfg.DB.DBURL = env.String("DB_URL", "")
	cfg.DB.SQLitePath = env.String("DB_SQLITE_PATH", "")

	cfg.DB.MaxIdleConns = env.Int("DB_MAX_IDLE_CONNS", 10)
	cfg.DB.MaxOpenConns = env.Int("DB_MAX_OPEN_CONNS", 100)
	cfg.DB.ConnMaxLifetime = env.Duration("DB_CONN_MAX_LIFETIME", time.Hour)
	cfg.DB.AutoMigrate = env.Bool("DB_AUTOMIGRATE", false)

	cfg.CountViews = countviews.New()
	featuresManager := features.NewManager(cfg.CountViews)
	if err := featuresManager.LoadAndValidate(); err != nil {
		return Config{}, err
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	switch c.DB.Driver {
	case "postgres":
		if c.DB.DBURL == "" {
			return errors.New("DB_URL must be set for postgres")
		}
	case "sqlite":
		if c.DB.SQLitePath == "" {
			return errors.New("DB_SQLITE_PATH must be set for sqlite")
		}
	default:
		return fmt.Errorf("unsupported DB_DRIVER: %q", c.DB.Driver)
	}
	if strings.TrimSpace(c.Server.Port) == "" {
		return errors.New("PORT must not be empty")
	}
	return nil
}
