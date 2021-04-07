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

	insert,err := db.Query("INSERT INTO users VALUE ('Elliot')")
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()

	fmt.Println("successfully connected")
	//http.HandleFunc("/hello",hello)
	//http.HandleFunc("/headers",headers)
	//http.ListenAndServe(":8090",nil)
}