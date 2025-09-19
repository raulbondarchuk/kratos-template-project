package ensure

import (
	"embed"
	"fmt"

	"gorm.io/gorm"
)

//go:embed ensure_schema.sql
var ensureFS embed.FS

func EnsureSchema(DB *gorm.DB, name string) error {
	// Read SQL template
	scriptContent, err := ensureFS.ReadFile("ensure_schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read ensure schema script: %v", err)
	}

	// Name of the database
	sqlScript := fmt.Sprintf(string(scriptContent), name)

	// Get raw connection
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %v", err)
	}

	// Check if schema exists
	var exists string
	err = sqlDB.QueryRow(sqlScript).Scan(&exists)
	if err != nil {
		panic(fmt.Sprintf("Error checking schema: %v", err))
	}

	if exists == "" {
		panic(fmt.Sprintf("Schema '%s' does not exist. Please create it manually before running the application", name))
	}

	return nil
}
