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

		queryMap := req.URL.Query()
		if queryMap ==nil {
			return
		}
		Username := queryMap["Username"]
		FirstName := queryMap["FirstName"]
		LastName := queryMap["LastName"]
		Email := queryMap["Email"]
		Password := queryMap["Password"]
		Phone := queryMap["Phone"]
		UserStatus := queryMap["UserStatus"]

		next_id := id + 1

		db.Query("INSERT into users value(?,?,?,?,?,?,?,?)",
			Username[0],FirstName[0],LastName[0],Email[0],Password[0],Phone[0],UserStatus[0],next_id)
		fmt.Println(queryMap)
	}
}
//func update_user(w http.ResponseWriter,req *http.Request) {
//	db,err := open_db()
//	if err != nil {
//		panic(err.Error())
//	}
//	defer db.Close()
//
//	//transactionの開始
//	tx,_ := db.Begin()
//	//最後にロールバック
//	defer tx.Rollback()
//	//auto incrementで追加
//	rows,_ := tx.Query("SELECT max(id) FROM users")
//	for rows.Next() {
//		var id int
//		rows.Scan(&id)
//		fmt.Println(id)
//
//		queryMap := req.URL.Query()
//		if queryMap ==nil {
//			return
//		}
//		Username := queryMap["Username"]
//		FirstName := queryMap["FirstName"]
//		LastName := queryMap["LastName"]
//		Email := queryMap["Email"]
//		Password := queryMap["Password"]
//		Phone := queryMap["Phone"]
//		UserStatus := queryMap["UserStatus"]
//
//		next_id := id + 1
//
//		db.Query("Update users SET value(?,?,?,?,?,?,?) where id = ?",
//			Username[0],FirstName[0],LastName[0],Email[0],Password[0],Phone[0],UserStatus[0],next_id)
//		fmt.Println(queryMap)
//	}
//}

func delete_user(w http.ResponseWriter, req *http.Request) {
	db,err := open_db()
	if err != nil{
		return
	}
	defer db.Close()
	tx,_ := db.Begin()
	defer tx.Rollback()

	queryMap := req.URL.Query()
	if queryMap ==nil {
		return
	}
	deleteUserId := queryMap["user_id"][0]
	db.Query("DELETE FROM users where id = ?", deleteUserId)
	fmt.Println(deleteUserId)
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
	http.HandleFunc("/create_user/",create_user)
	//http.HandleFunc("/update_user/",update_user)
	http.HandleFunc("/delete_user/",delete_user)
	http.HandleFunc("/headers",headers)
	http.ListenAndServe(":8090",nil)
}