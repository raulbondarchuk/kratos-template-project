package mysql

import (
	"fmt"

	"service/internal/conf/v1"
	"service/internal/data/adapters"
	"service/pkg/utils"
)

type adapter struct{}

func init() { adapters.Register("mysql", adapter{}) }

func (adapter) Name() string { return "mysql" }

func (adapter) LoadConfig(c *conf.Data, withSchema bool) (source string, logDSN string) {

	user := utils.EnvFirst("DB_USER")
	pass := utils.EnvFirst("DB_PASSWORD")
	host := utils.EnvFirst("DB_HOST")
	port := utils.EnvFirst("DB_PORT")
	db := utils.EnvFirst("DB_SCHEMA")

	if withSchema {
		source = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True&loc=Local",
			user, pass, host, port, db)
		logDSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True&loc=Local",
			user, "<password>", host, port, db)
	} else {
		source = fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=True&loc=Local",
			user, pass, host, port)
		logDSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=True&loc=Local",
			user, "<password>", host, port)
	}
	return source, logDSN
}
