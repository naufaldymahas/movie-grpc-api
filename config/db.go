package config

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var dbConn *sql.DB

func initDB() {
	if dbConn == nil {
		user := GetStringEnv("MYSQL_DB_USER", "root")
		password := GetStringEnv("MYSQL_DB_PASSWORD", "password")
		dbName := GetStringEnv("MYSQL_DB_NAME", "")

		dsn := fmt.Sprintf("%s:%s@/%s", user, password, dbName)

		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatal("Failed to initialized database connection")
		}
		db.SetConnMaxLifetime(time.Minute * 3)
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)

		dbConn = db
	}
}

func GetConn() *sql.DB {
	if dbConn == nil {
		initDB()
	}
	return dbConn
}
