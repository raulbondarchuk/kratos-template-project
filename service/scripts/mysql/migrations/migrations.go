package migrations

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

// AutoMigrate execute database migrations
func AutoMigrate(db *gorm.DB, models []any, logger log.Logger) error {
	h := log.NewHelper(logger)
	h.Info("Running database migrations...")

	// FK only works in InnoDB
	tx := db.Session(&gorm.Session{}).Set("gorm:table_options", "ENGINE=InnoDB")

	for _, m := range models {
		// hint in logs, what is being created right now
		h.Infof("migrating: %T", m)
		if err := tx.AutoMigrate(m); err != nil {
			h.Errorf("failed on %T: %v", m, err)
			return err
		}
		// sanity-check: table exists
		if !tx.Migrator().HasTable(m) {
			return fmt.Errorf("table for %T not found right after migration", m)
		}
	}

	h.Info("Database migration completed successfully")
	return nil
}
