package main

import (
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

func TestGetUser(t *testing.T)  {
	tt := []struct{
		getXToken string
		wantXToken string
		wantUserID int64
	}{
		{getXToken: "BpLnfgDsc2WD8F2qNfHK",wantXToken:"BpLnfgDsc2WD8F2qNfHK",wantUserID:1},
		{getXToken: "5a84jjJkwzDkh9h2fhfU",wantXToken:"5a84jjJkwzDkh9h2fhfU",wantUserID:2},
		{getXToken: "notexistXtoken",wantXToken:"",wantUserID:0},
	}

	for _ , tc := range tt{
		//arrange
		url :=  "localhost:8090/user/get/"
		req,err := http.NewRequest("GET",url,nil)
		//xtokenを設定
		req.Header.Set("xToken",tc.getXToken)

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

func TestFetchGacha(t *testing.T) {

	tt := []struct{
		getXToken string
		wantXToken string
		getTimes     string
		expectTimes string
	}{
		{getXToken: "BpLnfgDsc2WD8F2qNfHK",wantXToken:"BpLnfgDsc2WD8F2qNfHK",getTimes:"1" ,expectTimes:"1"},
		{getXToken: "BpLnfgDsc2WD8F2qNfHK",wantXToken:"BpLnfgDsc2WD8F2qNfHK",getTimes:"2" ,expectTimes:"2"},
		{getXToken: "5a84jjJkwzDkh9h2fhfU",wantXToken:"5a84jjJkwzDkh9h2fhfU",getTimes:"1" ,expectTimes:"1"},
		{getXToken: "5a84jjJkwzDkh9h2fhfU",wantXToken:"5a84jjJkwzDkh9h2fhfU",getTimes:"2" ,expectTimes:"2"},
		{getXToken: "5a84jjJkwzDkh9h2fhfU",wantXToken:"5a84jjJkwzDkh9h2fhfU",getTimes:"100" ,expectTimes:"2"},

	}

	for _ , tc := range tt{
		req,err := http.NewRequest("GET","localhost:8090/gacha/draw/?times="+tc.getTimes,nil)
		if err != nil{
			t.Fatalf("could not create request %v",err)
		}
		req.Header.Set("xToken",tc.getXToken)

		rec := httptest.NewRecorder()
		drawGacha(rec,req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK:got %v",res.Status)
		}
		data, err := ioutil.ReadAll(rec.Body)

		_, err = strconv.Atoi(string((bytes.TrimSpace(data))))

		gacha :=  make(map[string]interface{})

		if err := json.Unmarshal(data, &gacha); err != nil {
			log.Fatal(err)
		}
		arr := gacha["data"].([]interface{})

		assert.Equal(t,
			tc.expectTimes,
			strconv.Itoa(len(arr)),
			"Fetched Gacha data count is "+tc.getTimes,
		)
	}
}

type userTest struct {
	Name string
	Firstname string
	Lastname string
	Email string
	Password string
	Phone string
}


//
//func TestCreateUser(t *testing.T) {
//
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
//		req,err := http.NewRequest("GET","localhost:8090/create/?times="+tc.getUser,nil)
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
//}