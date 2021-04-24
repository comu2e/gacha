package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

// データベースに接続する
func DbInit() (*sql.DB, error) {
	var err error
	Db, err = sql.Open("mysql", "root:password@/testdb")

	return Db, err
}

// データベースから切断する
func DbClose() {
	if Db != nil {
		Db.Close()
	}
}

// データベースハンドラを取得する
func DbConn() *sql.DB {
	return Db
}
