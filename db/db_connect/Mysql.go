package db_connect

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"

	//_ "github.com/go-sql_nil-driver/mysql"
	"github.com/jmoiron/sqlx"
)
type MysqlConnectConfiguration struct {
	Host     string
	Database string
	User     string
	Password string
	Port     string
}

func CreateMysqlDBConnection(dbconf MysqlConnectConfiguration) (*sql.DB, error) {
	dsn := "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	dsn = fmt.Sprintf(dsn, dbconf.User, dbconf.Password, dbconf.Host, dbconf.Port, dbconf.Database)
	db, err := sql.Open("mysql", dsn)

	if err = db.Ping(); err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(10)
	go checkPing(db, time.Second*15)
	return db, nil
}

func CreateSQLXConnection(conn *sql.DB, driverName string) (db *sqlx.DB,err error) {
	db = sqlx.NewDb(conn, driverName)
	err = db.Ping()
	return
}

func CreateGormMysqlDBConnection(mysqlDB *sql.DB) (db *gorm.DB,err error) {
	db, err = gorm.Open(mysql.New(mysql.Config{
		Conn: mysqlDB,
	}), &gorm.Config{})
	return
}


func checkPing(DB *sql.DB, d time.Duration) {
	for {
		time.Sleep(d)
		err := DB.Ping()
		if err != nil {
			log.Println("error when checkPing", err)
		}
	}
}
