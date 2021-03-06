package main

import (
	"Gacha/database"
	"Gacha/model"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type userTest struct {
	Name string
	Firstname string
	Lastname string
	Email string
	Password string
	Phone string
}
const dbName = "root:password@/testdb"

func TestFetchXtoken(t *testing.T) {
	_, err := database.DbInit(dbName)
	if recover();err != nil {
		panic(err)
	}
	defer database.DbClose()

	tt := []struct{
		giveUsername string
		givePassword string
		wantToken string
	}{
		{giveUsername: "test",givePassword: "password",wantToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.P4Lqll22jQQJ1eMJikvNg5HKG-cKB0hUZA9BZFIG7Jk"},
	}

	for _,tc := range tt {
		//arrange
		query := "?Name="+tc.giveUsername+"&Password="+tc.givePassword
		url :=  "localhost:8090/user/fetch/"
		req, _ := http.NewRequest("GET",url+query,nil)
		//act
		rec := httptest.NewRecorder()
		fetchXtoken(rec,req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK:got %v",res.Status)
		}
		data,  _ := ioutil.ReadAll(rec.Body)

		_, err = strconv.Atoi(string(bytes.TrimSpace(data)))
		gacha :=  make(map[string]interface{})

		if err := json.Unmarshal(data, &gacha); err != nil {
			log.Fatal(err)
		}
		//assertion
		//token が得られているか確認する
		assert.Equal(t,
			tc.wantToken,
			gacha["data"],
		)
	}
}
func TestGetUser(t *testing.T)  {
	_, err := database.DbInit(dbName)
	if err != nil {
		panic(err)
	}
	defer database.DbClose()
	tt := []struct{
		giveXToken string
		wantXToken string
		wantUserID int64
	}{
		{giveXToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.efjmScJd31_IesdVdNsnd0i1jHE9rqAi28PXOMeSWLI",wantXToken:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.efjmScJd31_IesdVdNsnd0i1jHE9rqAi28PXOMeSWLI",wantUserID: 1},
		{giveXToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.P4Lqll22jQQJ1eMJikvNg5HKG-cKB0hUZA9BZFIG7Jk",wantXToken:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.P4Lqll22jQQJ1eMJikvNg5HKG-cKB0hUZA9BZFIG7Jk",wantUserID: 2},
	}

	for _ , tc := range tt{
		//arrange
		url :=  "localhost:8090/user/get/"
		req,err := http.NewRequest("GET",url,nil)
		//xtokenを設定
		req.Header.Set("xToken",tc.giveXToken)

		if err != nil{
			t.Fatalf("could not create request %v",err)
		}

		//act
		rec := httptest.NewRecorder()
		getUser(rec,req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK:got %v",res.Status)
		}
		data, err := ioutil.ReadAll(rec.Body)

		_, err = strconv.Atoi(string(bytes.TrimSpace(data)))

		var userJson model.UserJson
		if err := json.Unmarshal(data, &userJson); err != nil {
			log.Fatal(err)
		}
		//assertion
		assert.Equal(t,
			tc.wantUserID,
			userJson.Data.ID,
			"Fetched UserID is %v ",userJson.Data.ID,
		)

	}

}
//Create user

//
//func TestCreateUser(t *testing.T) {
//	_, err := database.DbInit()
//	if recover();err != nil {
//		panic(err)
//	}
//	defer database.DbClose()
//	tt := []struct {
//		inputUser userTest
//		wanUser userTest
//	}{
//		{inputUser: userTest{
//			Name:      "",
//			Firstname: "",
//			Lastname:  "",
//			Email:     "",
//			Password:  "",
//			Phone:     "",
//		},
//		wanUser: userTest{
//			Name:      "",
//			Firstname: "",
//			Lastname:  "",
//			Email:     "",
//			Password:  "",
//			Phone:     "",
//		},
//		},
//	}
//
//
//	for  _,tc := range tt{
//		length_userForm := 6
//		query := ""
//		for i := 0; i < length_userForm; i++ {
//			query += tc.inputUser.
//		}
//		req,err := http.NewRequest("GET","localhost:8090/create/?times="+tc.,nil)
//		if err != nil{
//			t.Fatalf("could not create request %v",err)
//		}
//
//		rec := httptest.NewRecorder()
//		drawGacha(rec,req)
//
//		res := rec.Result()
//		defer res.Body.Close()
//
//		if res.StatusCode != http.StatusOK {
//			t.Errorf("Expected status OK:got %v",res.Status)
//		}
//		data, err := ioutil.ReadAll(rec.Body)
//
//		_, err = strconv.Atoi(string((bytes.TrimSpace(data))))
//
//		gacha :=  make(map[string]interface{})
//
//		if err := json.Unmarshal(data, &gacha); err != nil {
//			log.Fatal(err)
//		}
//		arr := gacha["data"].([]interface{})
//
//		assert.Equal(t,
//			tc.expectTimes,
//			strconv.Itoa(len(arr)),
//			"Fetched Gacha data count is "+tc.getTimes,
//		)
//	}
//

//Update user
//delete user
func TestDrawGacha(t *testing.T) {
	_, err := database.DbInit(dbName)
	if recover();err != nil {
		panic(err)
	}
	defer database.DbClose()
	tt := []struct{
		giveXToken string
		wantXToken string
		giveTimes     string
		expectTimes string
	}{
		{
			giveXToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.efjmScJd31_IesdVdNsnd0i1jHE9rqAi28PXOMeSWLI",
			wantXToken:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.efjmScJd31_IesdVdNsnd0i1jHE9rqAi28PXOMeSWLI",
			giveTimes:"1" ,expectTimes:"1"},
	}

	for _ , tc := range tt{
		req,err := http.NewRequest("GET","localhost:8090/gacha/draw/?times="+tc.giveTimes,nil)
		if err != nil{
			t.Fatalf("could not create request %v",err)
		}
		req.Header.Set("xToken",tc.giveXToken)

		rec := httptest.NewRecorder()
		drawGacha(rec,req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK:got %v",res.Status)
		}
		data, err := ioutil.ReadAll(rec.Body)

		_, err = strconv.Atoi(string(bytes.TrimSpace(data)))

		gacha :=  make(map[string]interface{})

		if err := json.Unmarshal(data, &gacha); err != nil {
			log.Fatal(err)
		}
		arr := gacha["data"].([]interface{})

		assert.Equal(t,
			tc.expectTimes,
			strconv.Itoa(len(arr)),
			"Fetched Gacha data count is "+tc.giveTimes,
		)
	}
}
//getCharacterList
