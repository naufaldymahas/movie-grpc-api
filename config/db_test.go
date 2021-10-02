package config

import (
	"os"
	"testing"
)

func setDBEnv() {
	os.Setenv("MYSQL_DB_USER", "root")
	os.Setenv("MYSQL_DB_PASSWORD", "password")
	os.Setenv("MYSQL_DB_NAME", "bibit-test")
}

func TestGetConn(t *testing.T) {
	setDBEnv()
	conn := GetConn()
	if conn == nil {
		t.Error("Failed connect to database")
	}
}

func TestGetConn_Ping(t *testing.T) {
	setDBEnv()
	conn := GetConn()

	err := conn.Ping()

	if err != nil {
		t.Errorf("Failed to ping database and get error %v", err)
	}

}
