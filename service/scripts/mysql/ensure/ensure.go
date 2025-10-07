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

	// Check if database exists
	var dbName string
	err = sqlDB.QueryRow(sqlScript).Scan(&dbName)
	if err != nil {
		panic(fmt.Sprintf("Database '%s' does not exist", name))
	}

	return nil
}
