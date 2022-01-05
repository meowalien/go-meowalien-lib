package db_connect

import (
	"database/sql"
	"fmt"
	"github.com/meowalien/go-meowalien-lib/db/config_modules"
	"log"
	"time"

	//_ "github.com/go-sql_nil-driver/mysql"
	"github.com/jmoiron/sqlx"
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
	go checkPing(db.DB , time.Second*15)
	return db, nil
}

func checkPing(DB *sql.DB,d time.Duration) {
	for {
		time.Sleep(d)
		err := DB.Ping()
		if err != nil {
			log.Println("error when checkPing",err)
		}
	}
}