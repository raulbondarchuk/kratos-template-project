package postgres

import (
	"fmt"
	"service/internal/data/migrations"
	"service/pkg/utils"
	"strings"
	"time"

	pgEnsure "service/scripts/postgres/ensure"
	pgMigs "service/scripts/postgres/migrations"
	pgSeed "service/scripts/postgres/seed"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

type gormKratosWriter struct{ h *log.Helper }

func (w gormKratosWriter) Printf(format string, args ...interface{}) { w.h.Infof(format, args...) }

func (adapter) Connect(dsn string, logger log.Logger) (*gorm.DB, error) {
	h := log.NewHelper(logger)
	gLogger := glogger.New(gormKratosWriter{h}, glogger.Config{
		SlowThreshold: 300 * time.Millisecond,
		LogLevel:      glogger.Warn,
		Colorful:      true,
	})
	return gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gLogger})
}

func (adapter) EnsureSchema(db *gorm.DB) error {
	// In our model DB_SCHEMA — is the name of the database.
	target := utils.EnvFirst("DB_SCHEMA")
	// CREATE DATABASE IF NOT EXISTS in PG is not there — ignore "already exists".
	res := db.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, target))
	if res.Error != nil && !strings.Contains(res.Error.Error(), "already exists") {
		return res.Error
	}
	// If you want to do the schema inside the DB — do it after connecting to the target DB.
	return pgEnsure.EnsureSchema(db, target) // skeleton in your package; can be no-op
}

func (adapter) RunMigrations(db *gorm.DB, logger log.Logger) error {
	return pgMigs.AutoMigrate(db, migrations.MODELS_TO_MIGRATE, logger)
}
func (adapter) RunSeeds(db *gorm.DB, logger log.Logger) error { return pgSeed.SeedDatabase(db, logger) }
