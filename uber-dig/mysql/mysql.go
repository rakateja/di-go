package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	config "github.com/rakateja/di-go/uber-dig/config"
)

func New(cfg config.Config) (*sql.DB, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", cfg.MysqlUser, cfg.MysqlPassword, cfg.MysqlHost, cfg.MysqlDB)
	sqlDB, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}
	return sqlDB, nil
}
