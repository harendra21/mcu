package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var nameStartsWith string
var offset int = 0

func main() {

	fmt.Println("Enter your charecter name or just hit enter to start: ")
	fmt.Scanf("%s", &nameStartsWith)
	fmt.Println("--------------------------------------------")
	results := get_marvel_data()
	print_results(results)

	nextAction()

}

func nextAction() {
	fmt.Println("****************************************")
	fmt.Println("Choose your next action : ")
	fmt.Println("Type `next` for next page")
	fmt.Println("Type `prev` for previous page")
	fmt.Println("Type `new` for new search")
	fmt.Println("Type `exit` to quit")
	fmt.Println("****************************************")
	var action string
	fmt.Scanf("%s", &action)

	switch action {
	case "next":
		offset = offset + 10
	case "prev":
		if offset > 0 {
			offset = offset - 10
		} else {
			fmt.Println("****************************************")
			fmt.Println("No previous page")
			fmt.Println("****************************************")
			nextAction()
		}
	case "new":
		fmt.Println("Enter your charecter name or just hit enter to start: ")
		fmt.Scanf("%s", &nameStartsWith)
		offset = 0
	case "exit":
		os.Exit(0)
	default:
		fmt.Println("****************************************")
		fmt.Println("Action not recognized, Please try again")
		fmt.Println("****************************************")
		nextAction()
	}

	results := get_marvel_data()
	print_results(results)
	nextAction()

}

func print_results(results []interface{}) {
	for _, v := range results {
		data_json, _ := json.Marshal(v)
		var results map[string]interface{}
		json.Unmarshal([]byte(data_json), &results)
		thumb := results["thumbnail"].(map[string]interface{})
		fmt.Println("Name: ", results["name"])
		fmt.Println("Description: ", results["description"])
		fmt.Println("Image: ", thumb["path"])
		fmt.Println("-------------------------------------")
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
}

func get_marvel_data() []interface{} {
	limit := 10
	marvel_public_key := "302a803057c0480f5186b710a5454fc5"
	marvel_private_key := "ae2a810470116365dfb19025b37a5e39df4c7a1e"
	tNow := time.Now()
	tUnix := tNow.Unix()
	hash_str := fmt.Sprintf("%d%s%s", tUnix, marvel_private_key, marvel_public_key)
	hash_data := []byte(hash_str)
	hash := md5.Sum(hash_data)

	url := fmt.Sprintf("https://gateway.marvel.com/v1/public/characters?ts=%d&apikey=%s&hash=%x&limit=%d&offset=%d", tUnix, marvel_public_key, hash, limit, offset)

	if nameStartsWith != "" {
		url = fmt.Sprintf("%s&nameStartsWith=%s", url, nameStartsWith)
	}

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
