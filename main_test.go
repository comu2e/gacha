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
	tt := []struct{
 	times string
	}{{times:"1"},{times:"2"},{times:"100"}}

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

		//var character model.Character

		gacha :=  make(map[string]interface{})

		if err := json.Unmarshal(data, &gacha); err != nil {
			log.Fatal(err)
		}
		arr := gacha["data"].([]interface{})

		assert.Equal(t, tc.times,strconv.Itoa(len(arr)),"Fetched Gacha data count is "+tc.times )

	}


}