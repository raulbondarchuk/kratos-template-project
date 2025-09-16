package adapters

import (
	"service/internal/conf/v1"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type Adapter interface {
	// Loads the config.Source (DSN) with/without schema and returns a safe DSN for logs.
	LoadConfig(c *conf.Data, withSchema bool) (source string, logDSN string)

	// Connects to the database of a specific driver.
	Connect(dsn string, logger log.Logger) (*gorm.DB, error)

	// Creates the database/schema (when withSchema=false).
	EnsureSchema(db *gorm.DB) error

	// Migrations/seeds for the specific driver.
	RunMigrations(db *gorm.DB, logger log.Logger) error
	RunSeeds(db *gorm.DB, logger log.Logger) error

	// Name of the driver for logs.
	Name() string
}

var registry = map[string]Adapter{}

func Register(name string, a Adapter) { registry[name] = a }

func Get(name string) (Adapter, bool) { a, ok := registry[name]; return a, ok }
