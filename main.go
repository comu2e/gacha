package main

import (
	"Gacha/database"
	"Gacha/model"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//func openDb() (*sql.DB, error) {
//	db, err := sql.Open("mysql", "root:password@/testdb")
//	return db, err
//}

func setHeader(w http.ResponseWriter,method string)http.ResponseWriter  {

	return w
}
func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			_, _ = fmt.Fprintf(w, "%v : %v\n", name, h)
		}
	}
}
func fetchXtoken(w http.ResponseWriter, req *http.Request) {
	defer setHeader(w,"GET")

	if req.Method == http.MethodGet {

		queryMap := req.URL.Query()
		if queryMap == nil {
			return
		}
		userName := queryMap["Name"][0]
		passWord := queryMap["Password"][0]

		querySQL := fmt.Sprintf("SELECT xToken from users where Name = \"%s\" and Password = \"%s\" LIMIT 1", userName, passWord)

		fmt.Println(querySQL)
		db := database.DbConn()
		rows, _ := db.Query(querySQL)
		defer rows.Close()

		for rows.Next() {
			var user model.User

			_ = rows.Scan(&user.XToken)

			fmt.Println(user.XToken)
			output := map[string]interface{}{
				"data":    user.XToken,
				"status":  true,
				"message": "user data is fetched",
			}
			fmt.Println(output)
			defer func() error {
				outjson, err := json.Marshal(output)
				if err != nil {
					return err
				}
				w.Header().Set("content-Type", "application/json")
				_, err = fmt.Fprint(w, string(outjson))
				return err
			}()

		}
	}
	return
}

func getUser(w http.ResponseWriter, req *http.Request) {
	defer setHeader(w,"GET")

	xToken := req.Header.Get("xToken")
	db := database.DbConn()

	defer db.Close()

	queryMap := req.URL.Query()
	if queryMap == nil {
		return
	}
	row_count,err := db.Query("SELECT COUNT(id) as userCount FROM users WHERE xToken=?",xToken)
	defer row_count.Close()

	type userCount struct {
		count int
	}

	if err != nil {
		panic(err)
	}

	var dbCount userCount

	for row_count.Next() {
		err = row_count.Scan(&dbCount.count)
	}

	if dbCount.count != 0 {
		rows, err := db.Query("SELECT * FROM users WHERE xToken=?",xToken)
		if err != nil {
			panic(err)
		}

		for rows.Next() {

			var user model.User
			err = rows.Scan(&user.ID, &user.Name,
				&user.Firstname, &user.Lastname,
				&user.Email, &user.Password,
				&user.Phone, &user.UserStatus,&user.XToken)


			output := map[string]interface{}{
				"data":    user,
				"message": "user data is fetched",
			}
			defer func() error {
				outjson, err := json.Marshal(output)
				if err != nil {
					return err
				}
				w.Header().Set("content-Type", "application/json")
				_, err = fmt.Fprint(w, string(outjson))
				return err
			}()

		}

	}else{
		output := map[string]interface{}{
			"data": model.User{
				ID:         0,
				Name:       "",
				Firstname:  "",
				Lastname:   "",
				Email:      "",
				Password:   "",
				Phone:      "",
				UserStatus: false,
				XToken:     "",
			},
			"message": "user is not exist.",
		}
		defer func() error {
			outjson, err := json.Marshal(output)
			if err != nil {
				return err
			}
			w.Header().Set("content-Type", "application/json")
			_, err = fmt.Fprint(w, string(outjson))
			return err
		}()
	}

	return
}
func createUser(w http.ResponseWriter, req *http.Request) {
	defer setHeader(w,"POST")

	if req.Method == http.MethodPost {
		db := database.DbConn()

		//transactionの開始
		tx, _ := db.Begin()

		//auto incrementで追加
		rows, err := db.Query("SELECT max(id) FROM users")
		//usernameをユニークにするためにusernameのリストを取得する。

		if err != nil {
			http.Error(w, err.Error(), 401)
			panic(err)
		}

		for rows.Next() {
			var id int
			err = rows.Scan(&id)
			queryMap := req.URL.Query()
			if queryMap == nil {
				http.Error(w, err.Error(), 401)
				panic(err)
			}

			valueQuery := ""
			columnQuery := ""
			for k, v := range queryMap {
				if k == "Username" {
					queryUsername := v[0]

					rowsCount, _ := db.Query("SELECT count(Name) as hasUserCreated  from users where Name = ?", queryUsername)

					for rowsCount.Next() {
						var hasUserCreated int
						err = rowsCount.Scan(&hasUserCreated)

						if hasUserCreated != 0 {
							//userがunique出ないときにjsonでstatus:falseを返す
							panic(err)
						}
					}
				}
				valueQuery += "\"" + v[0] + "\"" + ","
				columnQuery += k + ","
			}

			valueQuery += strconv.Itoa(id+1) + ","
			columnQuery += "id" + ","

			xToken := generateXToken(20)
			valueQuery += "\"" + xToken + "\""
			columnQuery += "xToken" + ","

			valueQuery = strings.TrimRight(valueQuery, ",")
			columnQuery = strings.TrimRight(columnQuery, ",")

			query := "(" + columnQuery + ") " + "VALUES (" + valueQuery + ");"
			fmt.Println(query)
			_, err := tx.Query("INSERT INTO users" + query)

			if err != nil {
				//	//失敗したらロールバック
				_ = tx.Rollback()
				http.Error(w, err.Error(), 401)
				return
			}
			////成功したらCommit
			_ = tx.Commit()
			output := map[string]interface{}{
				"x-token": xToken,
				"message": "The user account was successfully created.",
				"status":  true,
			}
			defer func() error {
				outjson, err := json.Marshal(output)
				if err != nil {
					return err
				}
				w.Header().Set("content-Type", "application/json")
				_, err = fmt.Fprint(w, string(outjson))
				return err
			}()
		}
	}
}
func generateXToken(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
func updateUser(w http.ResponseWriter, req *http.Request) {
	setHeader(w,"PUT")
	if req.Method == http.MethodPut {
		db := database.DbConn()

		//transactionの開始
		tx, err := db.Begin()
		//auto incrementで追加
		queryMap := req.URL.Query()
		if queryMap == nil {
			return
		}
		xToken := req.Header.Get("xToken")
		setQuery := ""
		for k, v := range queryMap {
			if k != "id" {
				setQuery += k + " = \"" + v[0] + "\"" + ","
			}
		}
		setQuery = strings.TrimRight(setQuery, ",")
		fmt.Println(setQuery)
		query := "UPDATE users SET " + setQuery + " WHERE xToken = " + xToken

		tx.Query(query)
		if err != nil {
			err = tx.Rollback()
			return
		}

		_ = tx.Commit()
	}
	return
}
func deleteUser(w http.ResponseWriter, req *http.Request) {
	setHeader(w,"DELETE")
	if req.Method == http.MethodDelete {
		db := database.DbConn()
		tx, _ := db.Begin()
		xToken := req.Header.Get("xToken")
		_, errDel := tx.Query("DELETE FROM users where xToken = ?", xToken)
		if errDel != nil {
			_ = tx.Rollback()
		}
		_ = tx.Commit()
	}

}
func drawGacha(w http.ResponseWriter,req *http.Request) {

	setHeader(w,"GET")
	if req.Method == http.MethodGet {
		db := database.DbConn()
		//transactionの開始
		//tx, _ := db.Begin()

		xToken := req.Header.Get("xToken")

		defer db.Close()

		queryMap := req.URL.Query()
		if queryMap == nil {
			return
		}
		drawTimes := queryMap["times"][0]
		rows, err := db.Query("SELECT id FROM users where xToken = ?",xToken)
		if err != nil {
			return
		}
		defer rows.Close()

		var user_id string
		for rows.Next() {
			err = rows.Scan(&user_id)
		}
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		rows, err = db.Query("SELECT name,id FROM characters ORDER BY RAND() LIMIT ?", drawTimes)
		if err != nil {
			return
		}
		defer rows.Close()

		characters := []model.Character{}
		insertQuery := "INSERT INTO user_character (user_id,character_id) VALUES "
		for rows.Next() {
			var character model.Character
			err = rows.Scan(&character.Name, &character.ID)
			characters = append(characters, character)
			insertQuery += "("+user_id +"," + strconv.FormatInt(character.ID,10)+"),"
		}
		insertQuery = strings.TrimRight(insertQuery, ",")
		insertQuery += ";"
		db.Query(insertQuery)

		output := map[string]interface{}{
			"data":    characters,
			"message": "character data is fetched",
		}
		defer func() error {
			outjson, err := json.Marshal(output)
			if err != nil {
				return err
			}
			w.Header().Set("content-Type", "application/json")
			_, err = fmt.Fprint(w, string(outjson))
			return err
		}()
	}
}
func getCharacterList(w http.ResponseWriter, res *http.Request) {
	setHeader(w,"GET")

	if res.Method == http.MethodGet {

		db := database.DbConn()

		defer db.Close()
		queryMap := res.URL.Query()
		if queryMap == nil {
			return
		}
		user_id := queryMap["user_id"][0]
		rows, err := db.Query("SELECT character_id FROM user_character where user_id = ?", user_id)
		if err != nil {
			panic(err)
		}
		characters := []model.Character{}
		for rows.Next() {
			var character model.Character
			err = rows.Scan(&character.ID)
			characters = append(characters, character)
		}
		output := map[string]interface{}{
			"data":    characters,
			"message": "character data",
		}
		defer func() error {
			outjson, err := json.Marshal(output)
			if err != nil {
				return err
			}
			w.Header().Set("content-Type", "application/json")
			_, err = fmt.Fprint(w, string(outjson))
			return err
		}()
	}
}

func aboutMethodMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()

		t2 := time.Now()
		next.ServeHTTP(w, r)

		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))

	}
}
func setHeaderMiddleWare(next http.HandlerFunc,method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", method)


	}
}

func main() {
	//connection pool
	_, err := database.DbInit()
	if err != nil {
		panic(err)
	}
	defer database.DbClose()
	mux := http.NewServeMux()

	fmt.Println("successfully Launched")
	mux.HandleFunc("/user/get/",   setHeaderMiddleWare(aboutMethodMiddleWare(getUser),"GET)"))
	mux.HandleFunc("/user/fetch/", aboutMethodMiddleWare(fetchXtoken))
	mux.HandleFunc("/user/create/", aboutMethodMiddleWare(createUser))
	mux.HandleFunc("/user/update/", aboutMethodMiddleWare(updateUser))
	mux.HandleFunc("/user/delete/", aboutMethodMiddleWare(deleteUser))
	mux.HandleFunc("/gacha/draw/", aboutMethodMiddleWare(drawGacha))
	mux.HandleFunc("/character/list/",aboutMethodMiddleWare(getCharacterList))
	mux.HandleFunc("/headers", headers)
	if err := http.ListenAndServe(":8090", mux); err !=nil{
		log.Fatal(err)
	}

}
