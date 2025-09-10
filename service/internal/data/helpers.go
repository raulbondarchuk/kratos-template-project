package data

import (
	"context"
	"fmt"
	"service/internal/conf"
	"service/pkg/utils"
	"service/scripts/mysql/ensure"
	"service/scripts/mysql/migrations"
	"service/scripts/mysql/seed"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

type gormKratosWriter struct{ h *log.Helper }

func (w gormKratosWriter) Printf(format string, args ...interface{}) {
	w.h.Infof(format, args...)
}

// buildDSNWithSchema creates a DSN with schema name
func buildDSNWithSchema() string {
	user := utils.EnvFirst("DB_USER")
	password := utils.EnvFirst("DB_PASSWORD")
	host := utils.EnvFirst("DB_HOST")
	port := utils.EnvFirst("DB_PORT")
	schema := utils.EnvFirst("DB_SCHEMA")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True&loc=Local",
		user, password, host, port, schema)
}

// buildDSNWithoutSchema creates a DSN without schema name for initial connection
func buildDSNWithoutSchema() string {
	user := utils.EnvFirst("DB_USER")
	password := utils.EnvFirst("DB_PASSWORD")
	host := utils.EnvFirst("DB_HOST")
	port := utils.EnvFirst("DB_PORT")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=True&loc=Local",
		user, password, host, port)
}

// buildDSNForLogging creates a safe-to-log database connection string
func buildDSNForLogging() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True&loc=Local",
		utils.EnvFirst("DB_USER"), "<password>", utils.EnvFirst("DB_HOST"),
		utils.EnvFirst("DB_PORT"), utils.EnvFirst("DB_SCHEMA"))
}

// loadDatabaseConfig loads database configuration from environment
func loadDatabaseConfig(config *conf.Data, h *log.Helper, withSchema bool) {
	config.Database.Driver = utils.EnvFirst("DB_DRIVER")
	if withSchema {
		config.Database.Source = buildDSNWithSchema()
	} else {
		config.Database.Source = buildDSNWithoutSchema()
	}
	h.Infof("DSN: %s", buildDSNForLogging())
}

// connectDatabase establishes a database connection with the given configuration
func connectDatabase(dsn string, klog log.Logger) (*gorm.DB, error) {
	h := log.NewHelper(klog)

	gLogger := glogger.New(
		gormKratosWriter{h: h},
		glogger.Config{
			SlowThreshold:             300 * time.Millisecond,
			LogLevel:                  glogger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return db, nil
}

// verifyConnection checks if the database connection is working
func verifyConnection(db *gorm.DB, driver string, h *log.Helper) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	h.Infof("Connected to database: %s", driver)
	return nil
}

// setupSchema ensures database schema exists and is properly configured
func ensureSchema(db *gorm.DB) error {
	schemaName := utils.EnvFirst("DB_SCHEMA")
	if schemaName == "" {
		return fmt.Errorf("DB_SCHEMA environment variable is not set")
	}

	if err := ensure.EnsureSchema(db, schemaName); err != nil {
		return fmt.Errorf("failed to ensure schema exists: %v", err)
	}

	return nil
}

// Run migrations
func runMigrations(db *gorm.DB, logger log.Logger) error {
	return migrations.AutoMigrate(db, MODELS_TO_MIGRATE, logger)
}

func runSeeds(db *gorm.DB, logger log.Logger) error {
	return seed.SeedDatabase(db, logger)
}
