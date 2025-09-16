package postgres

import (
	"fmt"

	"service/internal/conf/v1"
	"service/internal/data/adapters"
	"service/pkg/utils"
)

type adapter struct{}

func init() { adapters.Register("postgres", adapter{}) }

func (adapter) Name() string { return "postgres" }

func (adapter) LoadConfig(c *conf.Data, withSchema bool) (source string, logDSN string) {

	user := utils.EnvFirst("DB_USER")
	pass := utils.EnvFirst("DB_PASSWORD")
	host := utils.EnvFirst("DB_HOST")
	port := utils.EnvFirst("DB_PORT")
	db := utils.EnvFirst("DB_SCHEMA")
	ssl := utils.EnvFirst("DB_SSLMODE")
	if ssl == "" {
		ssl = "disable"
	}
	tz := utils.EnvFirst("DB_TZ")
	if tz == "" {
		tz = "UTC"
	}

	if withSchema {
		source = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
			host, port, user, pass, db, ssl, tz)
		logDSN = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
			host, port, user, "<password>", db, ssl, tz)
	} else {
		// Connect to the system database "postgres" to create the target
		source = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s TimeZone=%s",
			host, port, user, pass, ssl, tz)
		logDSN = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s TimeZone=%s",
			host, port, user, "<password>", ssl, tz)
	}
	return source, logDSN
}
