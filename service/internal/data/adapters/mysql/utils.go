package mysql

import (
	"service/internal/data/migrations"
	"service/pkg/utils"
	"time"

	mysqlEnsure "service/scripts/mysql/ensure"
	mysqlMigs "service/scripts/mysql/migrations"
	mysqlSeed "service/scripts/mysql/seed"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

// Common things for GORM-logger and list of models
type gormKratosWriter struct{ h *log.Helper }

func (w gormKratosWriter) Printf(format string, args ...interface{}) { w.h.Infof(format, args...) }

func (adapter) Connect(dsn string, logger log.Logger) (*gorm.DB, error) {
	h := log.NewHelper(logger)
	gLogger := glogger.New(gormKratosWriter{h}, glogger.Config{
		SlowThreshold: 300 * time.Millisecond,
		LogLevel:      glogger.Warn,
		Colorful:      true,
	})
	return gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: gLogger})
}

func (adapter) EnsureSchema(db *gorm.DB) error {
	target := utils.EnvFirst("DB_SCHEMA")
	return mysqlEnsure.EnsureSchema(db, target)
}

func (adapter) RunMigrations(db *gorm.DB, logger log.Logger) error {
	return mysqlMigs.AutoMigrate(db, migrations.MODELS_TO_MIGRATE, logger)
}
func (adapter) RunSeeds(db *gorm.DB, logger log.Logger) error {
	return mysqlSeed.SeedDatabase(db, logger)
}
