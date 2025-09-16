package data

import (
	"context"
	"service/internal/conf/v1"
	"service/internal/data/adapters" // common registry
	_ "service/internal/data/adapters/mysql"
	_ "service/internal/data/adapters/postgres"
	"service/pkg/utils"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type Data struct{ DB *gorm.DB }

func NewData(config *conf.Data, logger log.Logger) (*Data, func(), error) {
	h := log.NewHelper(logger)

	// 1) Select adapter by DB_DRIVER / ENV (LoadConfig sets it)
	drv := utils.EnvFirst("DB_DRIVER")

	adapter, ok := adapters.Get(drv)
	if !ok {
		// Try to guess: mysql by default
		if a, ok2 := adapters.Get("mysql"); ok2 {
			adapter = a
		} else {
			return nil, nil, ErrUnknownDriver(drv)
		}
	}

	// 2) DSN without schema → connection to the base DB → ensure
	source, logDSN := adapter.LoadConfig(config, false)
	h.Infof("DSN (base): %s", logDSN)

	baseDB, err := adapter.Connect(source, logger)
	if err != nil {
		return nil, nil, err
	}

	if config.Database.EnsureSchema {
		if err := adapter.EnsureSchema(baseDB); err != nil {
			return nil, nil, err
		}
	} else {
		h.Infof("[DATABASE] [SKIPPED] Ensure schema is disabled")
	}

	// 3) DSN with schema → main connection
	source, logDSN = adapter.LoadConfig(config, true)
	h.Infof("DSN (target): %s", logDSN)

	db, err := adapter.Connect(source, logger)
	if err != nil {
		return nil, nil, err
	}

	// 4) Check connection (common)
	if err := verifyConnection(db, adapter.Name(), h); err != nil {
		return nil, nil, err
	}

	// 5) Migrations/seeds
	if config.Database.Migrations {
		if err := adapter.RunMigrations(db, logger); err != nil {
			return nil, nil, err
		}
	} else {
		h.Infof("[DATABASE] [SKIPPED] Migrations is disabled")
	}

	if config.Database.Seed {
		if err := adapter.RunSeeds(db, logger); err != nil {
			return nil, nil, err
		}
	} else {
		h.Infof("[DATABASE] [SKIPPED] Seed is disabled")
	}

	sqlDB, _ := db.DB()
	cleanup := func() { _ = sqlDB.Close() }

	return &Data{DB: db}, cleanup, nil
}

type driverError string

func (e driverError) Error() string   { return "unknown database driver: " + string(e) }
func ErrUnknownDriver(d string) error { return driverError(d) }

// Common check connection
func verifyConnection(db *gorm.DB, driver string, h *log.Helper) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return err
	}
	h.Infof("Connected to database: %s", driver)
	return nil
}
