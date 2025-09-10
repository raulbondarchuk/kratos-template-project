package data

import (
	"service/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

// Data represents the data layer with database connection
type Data struct {
	DB *gorm.DB
}

// NewData creates a new Data instance with database connection
func NewData(config *conf.Data, logger log.Logger) (*Data, func(), error) {
	h := log.NewHelper(logger)

	// Load configuration without schema for initial connection
	loadDatabaseConfig(config, h, false)

	// Connect to base database for schema setup
	baseDB, err := connectDatabase(config.Database.Source, logger)
	if err != nil {
		return nil, nil, err
	}

	// Ensure schema
	if config.Database.EnsureSchema {
		if err := ensureSchema(baseDB); err != nil {
			return nil, nil, err
		}
	}

	// Update configuration with schema for main connection
	loadDatabaseConfig(config, h, true)

	// Connect to specific schema
	db, err := connectDatabase(config.Database.Source, logger)
	if err != nil {
		return nil, nil, err
	}

	// Verify connection
	if err := verifyConnection(db, config.Database.Driver, h); err != nil {
		return nil, nil, err
	}

	// Create data instance
	d := &Data{DB: db}

	// Run migrations
	if config.Database.Migrations {
		if err := runMigrations(db, logger); err != nil {
			return nil, nil, err
		}
	}

	// Run seeds if needed
	if config.Database.Seed {
		if err := runSeeds(db, logger); err != nil {
			return nil, nil, err
		}
	}

	// Setup cleanup
	sqlDB, _ := db.DB()
	cleanup := func() {
		_ = sqlDB.Close()
	}

	return d, cleanup, nil
}
