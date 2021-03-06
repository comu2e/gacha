package main

import (
	"Gacha/database"
	"Gacha/model"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)


func fetchXtoken(w http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodGet {
		queryMap := req.URL.Query()
		if queryMap == nil {
			return
		}
		userName := queryMap["Name"][0]
		passWord := queryMap["Password"][0]

		querySQL := fmt.Sprintf("SELECT xToken FROM users WHERE Name = \"%s\"  AND Password = \"%s\" LIMIT 1", userName, passWord)
		db := database.DbConn()

		rows := db.QueryRow(querySQL)

		var user model.User
		_ = rows.Scan(&user.XToken)
		output := map[string]interface{}{
			"data":    user.XToken,
			"status":  true,
			"message": "user data is fetched",
		}
		defer func()  {
			outJson, err := json.Marshal(output)
			if err != nil {
				log.Println("Error:", err)
			}
			_, err = fmt.Fprint(w, string(outJson))
			log.Println("Error:", err)
		}()
	}
	return
}

func getUser(w http.ResponseWriter, req *http.Request ){
	xToken := req.Header.Get("xToken")
	if len(xToken) != 0{
		db := database.DbConn()
		rowCount := db.QueryRow("SELECT COUNT(id) as userCount FROM users WHERE xToken=?",xToken)

		type userCount struct {
			count int
		}

		var dbCount userCount
		err := rowCount.Scan(&dbCount.count)
		if err != nil{
			return
		}
		if dbCount.count != 0 {
			rows:= db.QueryRow("SELECT * FROM users WHERE xToken=?",xToken)
				var user model.User
			 	err := rows.Scan(&user.ID, &user.Name, &user.Firstname, &user.Lastname,
			 		&user.Email, &user.Password,
					&user.Phone, &user.UserStatus,&user.XToken)
			if err !=nil{
				return
			}
				output := map[string]interface{}{
					"data":    user,
					"message": "user data is fetched",
				}
				defer func()  {
					outJson, err := json.Marshal(output)
					if err != nil {
						log.Println("Error:", err)
					}
					_, err = fmt.Fprint(w, string(outJson))
					log.Println("Error:", err)
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
		defer func()  {
			outJson, err := json.Marshal(output)
			if err != nil {
				log.Println("Error:", err)
			}
			_, err = fmt.Fprint(w, string(outJson))
			log.Println("Error:", err)

		}()
	}

	return
}
func createUser(w http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {
		db := database.DbConn()
		//transaction?????????
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}
		//auto increment?????????
		rows := db.QueryRow("SELECT max(id) FROM users")
		//username?????????????????????????????????username??????????????????????????????

		var id int
		err = rows.Scan(&id)
		queryMap := req.URL.Query()
		if queryMap == nil {
			log.Println("Error:", err)
			return
		}
		valueQuery := ""
		columnQuery := ""
		username := ""
		var hasUserCreated int

		for k, v := range queryMap {
			if k == "Name" {
				queryUsername := v[0]
				username = queryUsername
				rowsCount := db.QueryRow("SELECT count(Name) as hasUserCreated  from users where Name = ?", queryUsername)
				err = rowsCount.Scan(&hasUserCreated)
			}
			valueQuery += "\"" + v[0] + "\"" + ","
			columnQuery += k + ","
		}
		if hasUserCreated == 0 {
			valueQuery += strconv.Itoa(id+1) + ","
			columnQuery += "id" + ","

			// ??????????????????
			xToken := jwt.New(jwt.SigningMethodHS256)
			//username???secret key?????????
			t, err := xToken.SignedString([]byte(username))

			valueQuery += "\"" + t + "\""
			columnQuery += "xToken" + ","

			valueQuery = strings.TrimRight(valueQuery, ",")
			columnQuery = strings.TrimRight(columnQuery, ",")
			query := "(" + columnQuery + ") " + "VALUES (" + valueQuery + ");"
			rows = tx.QueryRow("INSERT INTO users" + query)
			if err != nil {
				//	//?????????????????????????????????
				err = tx.Rollback()
				log.Println("Error:", err)
				return
			}
			//???????????????Commit
			err = tx.Commit()
			if err !=nil {
				http.Error(w, err.Error(), 401)
				return
			}
			output := map[string]interface{}{
				"x-token": t,
				"message": "The user account was successfully created.",
				"status":  true,
			}
			defer func() {
				outJson, err := json.Marshal(output)
				if err != nil {
					log.Println("Error:", err)
				}
				_, err = fmt.Fprint(w, string(outJson))
				log.Println("Error:", err)
			}()
		}else{
			_,err = fmt.Fprintf(w,"user has created")
			if err != nil{
				return
			}
		}
	}
	}
func updateUser(_ http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPut {
		db := database.DbConn()

		//transaction?????????
		tx, _ := db.Begin()
		//auto increment?????????
		queryMap := req.URL.Query()
		if queryMap == nil {
			log.Println("Error:Query is not exist." )
		}
		xToken := req.Header.Get("xToken")
		setQuery := ""
		for k, v := range queryMap {
			if k != "id" {
				setQuery += k + " = \"" + v[0] + "\"" + ","
			}
		}
		setQuery = strings.TrimRight(setQuery, ",")

		query := "UPDATE users SET " + setQuery + " WHERE xToken = " + xToken

		tx.QueryRow(query)

		_ = tx.Commit()
	}
	return
}
func deleteUser(_ http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodDelete {
		db := database.DbConn()
		tx,err := db.Begin()
		xToken := req.Header.Get("xToken")

		if len(xToken) != 0{
			tx.QueryRow("DELETE FROM users where xToken = ?", xToken)
			if err != nil {
				log.Println("Error:", err)
				err = tx.Rollback()
				if err != nil {
					log.Println("Error:", err)
				}
			}
			err = tx.Commit()
			if err != nil {
				log.Println("Error:", err)
			}
		} else {
			log.Println("Error:", err)
		}
	}
}
func drawGacha(w http.ResponseWriter,req *http.Request) {

	if req.Method == http.MethodGet {
		xToken := req.Header.Get("xToken")
		if len(xToken) != 0 {
			db := database.DbConn()
			//transaction?????????
			queryMap := req.URL.Query()
			if queryMap == nil {
				return
			}
			drawTimes := queryMap["times"][0]
			row  := db.QueryRow("SELECT id FROM users where xToken = ?", xToken)

			var userId string
			err := row.Scan(&userId)
			if err != nil {
				return
			}

			var characters []model.Character
			insertQuery := "INSERT INTO user_character (user_id,character_id) VALUES "

			rows,err := db.Query("SELECT name,id FROM characters ORDER BY RAND() LIMIT ?",drawTimes)
			if err != nil {
				return
			}
			for rows.Next() {
				var character model.Character
				err =  rows.Scan(&character.Name, &character.ID)
				if err != nil {
					return
				}
				characters = append(characters, character)
				insertQuery += "(" + userId + "," + strconv.FormatInt(character.ID, 10) + "),"
			}
			insertQuery = strings.TrimRight(insertQuery, ",")
			insertQuery += ";"
			db.QueryRow(insertQuery)

			output := map[string]interface{}{
				"data":    characters,
				"message": "character data is fetched",
			}
			defer func()  {

				outJson, err := json.Marshal(output)

				if err != nil {
					log.Println("Error: !!", err)
				}
				_, err = fmt.Fprint(w, string(outJson))
				log.Println("Error:", err)

			}()
		}else{
			_,err := fmt.Fprintf(w,"xToken is not setting")
			if err != nil {
				log.Println("Error:", err)
			}
		}
	}
}
func createCharacter(w http.ResponseWriter,req *http.Request) {
	password := req.Header.Get("xToken")
	if password  == "root"{
	if req.Method == http.MethodPut {
			db := database.DbConn()
			queryMap := req.URL.Query()
			if queryMap == nil {
				return
			}
			characterName := queryMap["characterName"][0]
			query := "INSERT INTO characters(name) value (\""+characterName+"\")"
			db.QueryRow(query)
			log.Print(query)
			output := map[string]interface{}{
				"data":   characterName ,
				"message": "character data is craeted",
			}
			defer func()  {

				outJson, err := json.Marshal(output)

				if err != nil {
					log.Println("Error: !!", err)
				}
				_, err = fmt.Fprint(w, string(outJson))
				log.Println("Error:", err)

			}()
		}
	}else {
		log.Print("Not allowed delete")

		fmt.Fprintln(w,"You are not allowed to create")
	}
}
func deleteCharacter(w http.ResponseWriter,req *http.Request) {
	password := req.Header.Get("xToken")
	if password  == "root"{
		if req.Method == http.MethodDelete {
			db := database.DbConn()
			queryMap := req.URL.Query()
			if queryMap == nil {
				return
			}
			id := queryMap["id"][0]
			//cascade???????????????
			db.QueryRow("DELETE FROM characters where id = ?", id)

			output := map[string]interface{}{
				"data":   id ,
				"message": "character data is deleted",
			}
			defer func()  {

				outJson, err := json.Marshal(output)

				if err != nil {
					log.Println("Error: !!", err)
				}
				_, err = fmt.Fprint(w, string(outJson))
				log.Println("Error:", err)

			}()
		}
	}else{
		log.Print("Not allowed delete")
	fmt.Fprintln(w,"You are not allowed to change")
	}

}
func getCharacterList(w http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodGet {

		db := database.DbConn()
		xToken := req.Header.Get("xToken")
		query := "SELECT id as userID FROM users WHERE xToken = " + "\""+xToken +"\""
		row := db.QueryRow(query)

		var userID string
		_ = row.Scan(&userID)
		query = "SELECT character_id,name,count(character_id) as character_count " +
			"FROM characters JOIN user_character uc on characters.id = uc.character_id WHERE user_id =" +userID +
			" GROUP BY characters.id order by character_id asc"
		rows,err := db.Query(query)

		var characterUsers []model.CharacterUser

		if err !=  nil {
			return
		}
		for rows.Next() {
			var characterUser model.CharacterUser
			err := rows.Scan(&characterUser.CharacterID,&characterUser.Name,&characterUser.Character_count)
			if err != nil {
				return
			}
			characterUsers = append(characterUsers,characterUser)
		}

		output := map[string]interface{}{
			"data":characterUsers,
			"message": "characters data",
		}
		defer func()  {
			outJson, err := json.Marshal(output)

			if err != nil {
				log.Println("Error:", err)
			}
			_, err = fmt.Fprint(w, string(outJson))
			log.Println("Error:", err)
		}()
	}
}

func RequestLog(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tStart := time.Now()
		next.ServeHTTP(w, r)
		tEnd := time.Now()

		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), tEnd.Sub(tStart))
	}
}
func setHeaderMiddleWare(next http.HandlerFunc,method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", method)
		next.ServeHTTP(w, r)

	}
}

func main(){
	const  dbName = "root:password@/maindb"

	_, err := database.DbInit(dbName)
	if recover();err != nil {
		panic(err)
	}
	defer database.DbClose()
	mux := http.NewServeMux()
	fmt.Println("successfully Launched")
	mux.HandleFunc("/user/get/",   RequestLog(setHeaderMiddleWare(getUser,     "GET)")))
	mux.HandleFunc("/user/fetch/", RequestLog(setHeaderMiddleWare(fetchXtoken, "GET")))
	mux.HandleFunc("/user/create/",RequestLog(setHeaderMiddleWare(createUser,  "POST")))
	mux.HandleFunc("/user/update/",RequestLog(setHeaderMiddleWare(updateUser,  "PUT")))
	mux.HandleFunc("/user/delete/",RequestLog(setHeaderMiddleWare(deleteUser,  "DELETE")))
	mux.HandleFunc("/gacha/draw/", RequestLog(setHeaderMiddleWare(drawGacha,    "GET")))
	mux.HandleFunc("/character/create/", RequestLog(setHeaderMiddleWare(createCharacter,    "PUT")))
	mux.HandleFunc("/character/delete/", RequestLog(setHeaderMiddleWare(deleteCharacter,    "DELETE")))
	mux.HandleFunc("/character/list/", RequestLog(setHeaderMiddleWare(getCharacterList, "GET")))
	if err := http.ListenAndServe(":8090", mux); err != nil {
		log.Fatal(err)
	}
}
