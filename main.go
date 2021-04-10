package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go_cyber/models"
	"net/http"
	"log"
	"strings"
)

func openDb() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:password@/testdb")
	return db,err
}

func fetchUser(_ http.ResponseWriter, req *http.Request) {
	db,err := openDb()
	if err != nil {
		return
	}
	defer db.Close()
	queryMap := req.URL.Query()
	if queryMap ==nil {
		return
	}
	//id := queryMap["id"][0]

	rows,err := db.Query("SELECT * FROM users where id = ?",2)

	if err != nil{
		log.Fatal(err)
	}
	//TODO userの情報を取得する
	//TODO jsonに出力する。
	for rows.Next() {
		user := &models.User{}
		rows.Scan(&user.Username)
		fmt.Println(user)
	}



}


func createUser(_ http.ResponseWriter,req *http.Request) {
	db,err := openDb()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	//transactionの開始
	tx,_ := db.Begin()

	//auto incrementで追加
	rows,err := db.Query("SELECT max(id) FROM users")

	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var id int
		_ = rows.Scan(&id)
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

		nextId := id + 1
		_, errInsert := tx.Query("INSERT into users value(?,?,?,?,?,?,?,?)",
			Username[0],FirstName[0],LastName[0],Email[0],Password[0],Phone[0],UserStatus[0], nextId)
		fmt.Println(queryMap)
		if errInsert != nil{
			//失敗したらロールバック
			_ = tx.Rollback()
			log.Fatal(errInsert)

		}
		//成功したらCommit
		_ = tx.Commit()
		}
	}

func updateUser(_ http.ResponseWriter,req *http.Request) {
	db,err := openDb()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	//transactionの開始
	var tx,_ = db.Begin()
	//auto incrementで追加
	queryMap := req.URL.Query()
		if queryMap ==nil {
			return
		}
	id := queryMap["id"][0]
	//TODO:UPDATE SQL 最大で7件発行されるので、１件にまとめられないか 予めquery文作成しておくことで対応。
	setQuery := ""
	for k,v := range queryMap{
		if k != "id"{
			setQuery += k + " = \"" + v[0] +"\"" +  ","
		}
	}
	setQuery = strings.TrimRight(setQuery, ",")
	fmt.Println(setQuery)
	query := "UPDATE users SET " + setQuery + " WHERE id = " + id
	fmt.Println(query)
	_, errUpdate := tx.Query(query)
	if errUpdate != nil {
		_ = tx.Rollback()
	}

	_ = tx.Commit()

}

func deleteUser(_ http.ResponseWriter, req *http.Request) {
	db,err := openDb()
	if err != nil{
		return
	}
	defer db.Close()
	tx,_ := db.Begin()

	queryMap := req.URL.Query()

	deleteUserId := queryMap["user_id"][0]
	_, errDel := tx.Query("DELETE FROM users where id = ?", deleteUserId)
	if errDel != nil{
		_ = tx.Rollback()
	}
	_ = tx.Commit()

}
func headers(w http.ResponseWriter, req *http.Request) {
	for name,headers := range req.Header{
		for _,h := range headers{
			_, _ = fmt.Fprintf(w, "%v : %v\n", name, h)
		}
	}
}

func main() {

	db,err := openDb()

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	fmt.Println("successfully connected")
	//TODO:GET,PUT,DELETEはどのように指定すればいいのか確認する
	http.HandleFunc("/user/get/", fetchUser)
	http.HandleFunc("/user/create/", createUser)
	http.HandleFunc("/user/update/", updateUser)
	http.HandleFunc("/user/delete/", deleteUser)
	http.HandleFunc("/headers",headers)
	_ = http.ListenAndServe(":8090", nil)
}