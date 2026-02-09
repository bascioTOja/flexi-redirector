package config

import (
	"errors"
	"fmt"
	"regexp"
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

	TableShortUrlName    string
	ColShortUrlId        string
	ColShortUrlName      string
	ColShortUrlSlug      string
	ColShortUrlLongUrl   string
	ColShortUrlViews     string
	ColShortUrlUpdatedAt string
	ColShortUrlCreatedAt string
	ColShortUrlCreatedBy string
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

	cfg.DB.TableShortUrlName = env.String("DB_TABLE_SHORT_URL_NAME", "short_urls")
	cfg.DB.ColShortUrlId = env.String("DB_COL_SHORT_URL_ID", "id")
	cfg.DB.ColShortUrlName = env.String("DB_COL_SHORT_URL_NAME", "name")
	cfg.DB.ColShortUrlSlug = env.String("DB_COL_SHORT_URL_SLUG", "slug")
	cfg.DB.ColShortUrlLongUrl = env.String("DB_COL_SHORT_URL_LONG_URL", "long_url")
	cfg.DB.ColShortUrlUpdatedAt = env.String("DB_COL_SHORT_URL_UPDATED_AT", "updated_at")
	cfg.DB.ColShortUrlCreatedAt = env.String("DB_COL_SHORT_URL_CREATED_AT", "created_at")
	cfg.DB.ColShortUrlCreatedBy = env.String("DB_COL_SHORT_URL_CREATED_BY", "created_by")

	// Load features
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

	requiredIdentifierFields := map[string]string{
		"DB_TABLE_SHORT_URL_NAME":   c.DB.TableShortUrlName,
		"DB_COL_SHORT_URL_ID":       c.DB.ColShortUrlId,
		"DB_COL_SHORT_URL_SLUG":     c.DB.ColShortUrlSlug,
		"DB_COL_SHORT_URL_LONG_URL": c.DB.ColShortUrlLongUrl,
	}
	for envName, identifierValue := range requiredIdentifierFields {
		if err := validateSQLIdentifier(envName, identifierValue); err != nil {
			return err
		}
	}

	// Optional columns: allow "-" or an empty value, which means "do not use".
	optionalIdentifierFields := map[string]string{
		"DB_COL_SHORT_URL_NAME":       c.DB.ColShortUrlName,
		"DB_COL_SHORT_URL_VIEWS":      c.DB.ColShortUrlViews,
		"DB_COL_SHORT_URL_UPDATED_AT": c.DB.ColShortUrlUpdatedAt,
		"DB_COL_SHORT_URL_CREATED_AT": c.DB.ColShortUrlCreatedAt,
		"DB_COL_SHORT_URL_CREATED_BY": c.DB.ColShortUrlCreatedBy,
	}
	for envName, identifierValue := range optionalIdentifierFields {
		if err := validateOptionalSQLIdentifier(envName, identifierValue); err != nil {
			return err
		}
	}

	return nil
}

var sqlIdentifierPattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

func validateSQLIdentifier(fieldName string, identifier string) error {
	trimmed := strings.TrimSpace(identifier)
	if trimmed == "" {
		return fmt.Errorf("%s must not be empty", fieldName)
	}
	if !sqlIdentifierPattern.MatchString(trimmed) {
		return fmt.Errorf("%s must be a valid SQL identifier (got %q)", fieldName, identifier)
	}
	return nil
}

func validateOptionalSQLIdentifier(fieldName string, identifier string) error {
	trimmed := strings.TrimSpace(identifier)
	if trimmed == "" || trimmed == "-" {
		return nil
	}
	return validateSQLIdentifier(fieldName, identifier)
}
