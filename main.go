package main

import (
	"Gacha/model"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
	"strings"
)

func openDb() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:password@/testdb")
	return db,err
}

func getUser(w http.ResponseWriter, req *http.Request)  {
	db,err := openDb()
	if err != nil {
		return
	}
	defer db.Close()
	queryMap := req.URL.Query()
	if queryMap ==nil {
		return
	}
	id := queryMap["id"][0]
	fmt.Println(id)

	rows,err := db.Query("SELECT * FROM users where id = ?",id)
	if err != nil{
		return
	}
	//TODO userの情報を取得する
	//TODO jsonに出力する。
	for rows.Next() {
		var user model.User

		err = rows.Scan(&user.Id,&user.Username,
			 			 &user.Firstname,&user.Lastname,
			 			 &user.Email,&user.Password,
			 			 &user.Phone,&user.UserStatus)
		fmt.Println(user.Id, user.Username)
		output := map[string]interface{}{
			//Todo id = 2で得られるが、id=2aとしても得られるので修正する。
			//Todo error のときのjsonも準備する。
			"data":user,
			"message":"user data is fetched",
		}
		defer func() error{
			outjson ,err := json.Marshal(output)
			if err != nil {
				return err
			}
			w.Header().Set("content-Type","application/json")
			_, err = fmt.Fprint(w, string(outjson))
			return err
		}()

	}

	return
}


func createUser(_ http.ResponseWriter,req *http.Request)  {
	if req.Method == http.MethodPost  {
		db,err := openDb()
		if err != nil {
			return
		}
		defer db.Close()

		//transactionの開始
		tx,_ := db.Begin()

		//auto incrementで追加
		rows,err := db.Query("SELECT max(id) FROM users")

		if err != nil {
			return
		}
		for rows.Next() {
			var id int
			err = rows.Scan(&id)
			fmt.Println(id)
			queryMap    := req.URL.Query()
			if queryMap == nil {
				return
			}

			valueQuery  := ""
			columnQuery := ""
			for k,v := range queryMap {
					valueQuery  += "\""+v[0] +"\"" + ","
					columnQuery += k + ","
			}

			valueQuery  += strconv.Itoa(id + 1) + ","
			columnQuery += "id" + ","
			valueQuery  = strings.TrimRight(valueQuery  , ",")
			columnQuery = strings.TrimRight(columnQuery, ",")

			fmt.Println(valueQuery)
			fmt.Println(columnQuery)

			query  := "("+columnQuery+") " + "VALUES("+valueQuery+");"
			_, err := tx.Query("INSERT into users"+ query)
			fmt.Println(query)
			if err != nil{
			//	//失敗したらロールバック
				_ = tx.Rollback()
				return
			//
			}
			////成功したらCommit
			_ = tx.Commit()
		}
	}

	return
}

func updateUser(w http.ResponseWriter,req *http.Request) {

	if  req.Method == http.MethodPut{
		db,err := openDb()
		if err != nil {
			return
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
		tx.Query(query)
		if err != nil {
			err = tx.Rollback()
			return
		}

		_ = tx.Commit()
	}
	return
}

func deleteUser(_ http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodDelete {

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
	http.HandleFunc("/user/get/", getUser)
	http.HandleFunc("/user/create/", createUser)
	http.HandleFunc("/user/update/", updateUser)
	http.HandleFunc("/user/delete/", deleteUser)
	http.HandleFunc("/headers",headers)
	_ = http.ListenAndServe(":8090", nil)
}