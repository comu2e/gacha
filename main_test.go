package main

import (
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

// 本来の処理のダミーその1



func TestMyHandler(t *testing.T) {
	tt := []struct {
	times         string
	expect_times string
}{
		{times:"1",expect_times:"1"},
		{times:"2",expect_times:"2"},
		//キャラクターが2つ登録されているときは2.登録数の上限が期待する値
		{times:"100",expect_times:"2"},
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
			tc.expect_times,
			strconv.Itoa(len(arr)),
			"Fetched Gacha data count is "+tc.times,
		)


	}


}