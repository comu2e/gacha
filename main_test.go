package main

import (
	"Gacha/model"
	"bytes"
	"encoding/json"
	"fmt"
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
		giveUserID string
		wantUserID string
	}{
		{giveUserID: "1",wantUserID: "1"},
		{giveUserID: "2",wantUserID: "2"},
	}

	for _ , tc := range tt{
		url :=  "localhost:8090/user/get/?id="+tc.giveUserID
		print(url+"\n")
		req,err := http.NewRequest("GET",url,nil)
		if err != nil{
			t.Fatalf("could not create request %v",err)
		}

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
		fmt.Println(userJson.Data.XToken)
	}

}

func TestFetchGacha(t *testing.T) {
	tt := []struct {
	times         string
	expectTimes string
}{
		{times:"1",expectTimes:"1"},
		{times:"2",expectTimes:"2"},
		//キャラクターが2つ登録されているときは2.登録数の上限が期待する値
		{times:"100",expectTimes:"2"},
	}

	for _ , tc := range tt{
		req,err := http.NewRequest("GET","localhost:8090/gacha/draw/?times="+tc.times,nil)
		if err != nil{
			t.Fatalf("could not create request %v",err)
		}

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
			"Fetched Gacha data count is "+tc.times,
		)
	}
}