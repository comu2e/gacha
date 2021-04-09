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

func createUser(w http.ResponseWriter,req *http.Request) {
	db,err := open_db()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	//transactionの開始
	tx,_ := db.Begin()

	//auto incrementで追加
	rows,err := db.Query("SELECT max(id) FROM users")

	if err != nil {
		 fmt.Println(err)
	}
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
		_,err_insert := tx.Query("INSERT into users value(?,?,?,?,?,?,?,?)",
			Username[0],FirstName[0],LastName[0],Email[0],Password[0],Phone[0],UserStatus[0],next_id)
		fmt.Println(queryMap)
		if err_insert != nil{
			//失敗したらロールバック
			tx.Rollback()
		}
		//成功したらCommit
		tx.Commit()
		}
	}

func updateUser(w http.ResponseWriter,req *http.Request) {
	db,err := open_db()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	//transactionの開始
	tx,err := db.Begin()

	//auto incrementで追加
	queryMap := req.URL.Query()
		if queryMap ==nil {
			return
		}
	id := queryMap["id"][0]
	//TODO:UPDATE SQL 最大で7件発行されるので、１件にまとめられないか
	for k, v := range queryMap {

		if k != "id" {
			//query文を作成
			query := "UPDATE users SET " + k + " = \"" + v[0] + "\" WHERE id = " + id
			fmt.Println(query)
			_,err_update := tx.Query(query)
			if err_update != nil {
				tx.Rollback()
			}
		}
	}
	tx.Commit()

}

func deleteUser(w http.ResponseWriter, req *http.Request) {
	db,err := open_db()
	if err != nil{
		return
	}
	defer db.Close()
	tx,_ := db.Begin()

	queryMap := req.URL.Query()

	deleteUserId := queryMap["user_id"][0]
	_,err_del := tx.Query("DELETE FROM users where id = ?", deleteUserId)
	if err_del != nil{
		tx.Rollback()
	}
	tx.Commit()

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
	//TODO:GET,PUT,DELETEはどのように指定すればいいのか確認する
	http.HandleFunc("/create_user/", createUser)
	http.HandleFunc("/update_user/", updateUser)
	http.HandleFunc("/delete_user/", deleteUser)
	http.HandleFunc("/headers",headers)
	http.ListenAndServe(":8090",nil)
}