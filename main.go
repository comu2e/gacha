package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

func open_db() (*sql.DB, error) {
	db, err := sql.Open("mysql", "mysql")
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
	if err != nil{
		//return nil error
		panic("cannot open database")
	}
	defer db.Close()
	insert, err := db.Query("INSERT INTO test VALUES ( 2, 'TEST' )")

	// if there is an error inserting, handle it
	if err != nil {
		panic(err.Error())
	}
	// be careful deferring Queries if you are using transactions
	defer insert.Close()

	http.HandleFunc("/hello",hello)
	http.HandleFunc("/headers",headers)
	http.ListenAndServe(":8090",nil)
}