package db

import (
	"fmt"

	"flexi-redirector/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Open(cfg config.DBConfig) (*gorm.DB, func() error, error) {
	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	var dialector gorm.Dialector
	switch cfg.Driver {
	case "postgres":
		dialector = postgres.Open(cfg.DBURL)
	case "sqlite":
		dialector = sqlite.Open(cfg.SQLitePath)
	default:
		return nil, nil, fmt.Errorf("unsupported driver: %s", cfg.Driver)
	}

	db, err := gorm.Open(dialector, gormCfg)
	if err != nil {
		return nil, nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	closeFn := func() error { return sqlDB.Close() }
	return db, closeFn, nil
}
