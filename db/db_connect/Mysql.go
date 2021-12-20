package db_connect

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/meowalien/go-meowalien-lib/db/config_modules"
)

func CreateMysqlDBConnectionWithSQLX(dbconf config_modules.MysqlConnectConfiguration) (*sqlx.DB, error) {
	dsn := "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	dsn = fmt.Sprintf(dsn, dbconf.User, dbconf.Password, dbconf.Host, dbconf.Port, dbconf.Database)
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(10)
	return db, nil
}
