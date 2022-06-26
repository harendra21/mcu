package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {

	results := get_marvel_data()

	for _, v := range results {
		data_json, _ := json.Marshal(v)
		var results map[string]interface{}
		json.Unmarshal([]byte(data_json), &results)
		thumb := results["thumbnail"].(map[string]interface{})
		fmt.Println(results["name"])
		fmt.Println(results["description"])
		fmt.Println(thumb["path"])
		fmt.Println("--------------------")
	}

}
func get_marvel_data() []interface{} {

	limit := 3

	marvel_public_key := "302a803057c0480f5186b710a5454fc5"
	marvel_private_key := "ae2a810470116365dfb19025b37a5e39df4c7a1e"
	tNow := time.Now()
	tUnix := tNow.Unix()
	hash_str := fmt.Sprintf("%d%s%s", tUnix, marvel_private_key, marvel_public_key)
	hash_data := []byte(hash_str)
	hash := md5.Sum(hash_data)

	url := fmt.Sprintf("https://gateway.marvel.com/v1/public/characters?ts=%d&apikey=%s&hash=%x&limit=%d", tUnix, marvel_public_key, hash, limit)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		panic(err)
	}
	data := result["data"]
	data_marshal, _ := json.Marshal(data)
	var results map[string]interface{}
	json.Unmarshal([]byte(data_marshal), &results)
	resp_results := results["results"].([]interface{})
	return resp_results

}
