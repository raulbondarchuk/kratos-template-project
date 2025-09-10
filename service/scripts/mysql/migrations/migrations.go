package migrations

import (
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

// AutoMigrate execute database migrations
func AutoMigrate(db *gorm.DB, models []any, logger log.Logger) error {
	h := log.NewHelper(logger)
	h.Info("Running database migrations...")

	// Add new models to the migration list
	err := db.AutoMigrate(
		models...,
	)
	if err != nil {
		h.Errorf("Failed to migrate database: %v", err)
		return err
	}

	h.Info("Database migration completed successfully")
	return nil
}
