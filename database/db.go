package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// データベースに接続する
func DbInit() (*sql.DB, error) {
	var err error
	db, err = sql.Open("mysql", "root:password@/testdb")

	return db, err
}

// データベースから切断する
func DbClose() {
	if db != nil {
		db.Close()
	}
}

// データベースハンドラを取得する
func DbConn() *sql.DB {
	return db
}
