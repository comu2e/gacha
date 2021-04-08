package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func open_db() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:password@/testdb")
	return db,err
}
func create_user(w http.ResponseWriter,req *http.Request) {
	db,err := open_db()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	//transactionの開始
	tx,_ := db.Begin()
	//最後にロールバック
	defer tx.Rollback()
	//auto incrementで追加
	rows,_ := tx.Query("SELECT max(id) FROM users")
	for rows.Next() {
		var id int
		rows.Scan(&id)
		fmt.Println(id)
		id_d := id + 1
		db.Query("INSERT into characters value(?,?,?,?,?,?,?,?)",
			Username,FirstName,LastName,Email,Password,Phone,UserStatus,id_d)
	}
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "CCCC\n")
}
func headers(w http.ResponseWriter, req *http.Request) {

	for name,headers := range req.Header{
		for _,h := range headers{
			fmt.Fprintf(w,"%v : %v\n",name,h)
		}
	}
}

func main() {
	db,err := open_db()

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	fmt.Println("successfully connected")
	http.HandleFunc("/create_user",create_user)
	http.HandleFunc("/headers",headers)
	http.ListenAndServe(":8090",nil)
}