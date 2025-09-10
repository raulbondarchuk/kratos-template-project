package seed

import (
	"embed"
	"fmt"

	"sort"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

//go:embed *.sql
var seedsFS embed.FS

// SeedDatabase executes all seed scripts
func SeedDatabase(DB *gorm.DB, logger log.Logger) error {
	h := log.NewHelper(logger)
	h.Info("Seeding database...")

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Read list of files
	entries, err := seedsFS.ReadDir(".")
	if err != nil {
		return fmt.Errorf("failed to read seeds dir: %w", err)
	}

	// To execute in order 01_, 02_, 03_
	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	// Execute all scripts
	for _, file := range files {
		script, err := seedsFS.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read seed %s: %w", file, err)
		}

		h.Infof("Seeding %s ...", file)
		if _, err := sqlDB.Exec(string(script)); err != nil {
			return fmt.Errorf("failed to execute seed %s: %w", file, err)
		}
	}

	return nil
}
