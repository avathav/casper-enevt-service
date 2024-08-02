package database

import (
	"database/sql"
	"fmt"
)

func mysqlConnectionString(params Parameters, config Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&timeout=%s&multiStatements=%t",
		params.Username,
		params.Password,
		params.Host,
		params.Port,
		params.Database,
		config.Charset,
		config.ParseTime,
		config.EstablishConnectionTimeout,
		config.MultiStatements,
	)
}

func MysqlConnection(params Parameters) (*sql.DB, error) {
	return sql.Open(DriverName, mysqlConnectionString(params, DefaultConnectionConfig()))
}
