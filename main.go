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

	"github.com/allegro/bigcache/v3"
)

var nameStartsWith string
var offset int = 0
var total int = 0
var limit int = 10
var cache *bigcache.BigCache
var testing_mode bool = false

type Response struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

type Data struct {
	Offset  int       `json:"offset"`
	Limit   int       `json:"limit"`
	Total   int       `json:"total"`
	Results []Results `json:"results"`
}

type Results struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Thumbnail   Thumbnail `json:"thumbnail"`
}
type Thumbnail struct {
	Path string `json:"path"`
}

func main() {
	cache, _ = bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))

	fmt.Println("Enter your charecter name or just hit enter to start: ")
	fmt.Scanf("%s", &nameStartsWith)
	fmt.Println("--------------------------------------------")
	results := getMarvelData()
	printResults(results)
	nextAction()
}

func nextAction() {
	fmt.Println("****************************************")
	fmt.Println("Choose your next action : ")
	fmt.Println("Type `next` for next page, `prev` for previous page, `new` for new search, `exit` to quit")
	fmt.Println("****************************************")
	var action string
	fmt.Scanf("%s", &action)

	switch action {
	case "next":
		offset = offset + limit
		if offset > total {
			fmt.Println("****************************************")
			fmt.Println("No more records")
			fmt.Println("****************************************")
			nextAction()
		}
	case "prev":
		if offset > 0 {
			offset = offset - limit
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

	results := getMarvelData()
	printResults(results)
	nextAction()

}

func printResults(results Data) {
	total = results.Total
	fmt.Println("Total: ", total)
	for _, v := range results.Results {
		thumb := v.Thumbnail
		fmt.Println("Name: ", v.Name)
		fmt.Println("Description: ", v.Description)
		fmt.Println("Image: ", thumb.Path)
		fmt.Println("------------------------------------------------------------------------")
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
}

func getMarvelData() Data {

	key := fmt.Sprintf("marvel_%s_%d", nameStartsWith, offset)
	entry, _ := cache.Get(key)
	if len(entry) > 0 && !testing_mode {
		var res Response
		if err := json.Unmarshal(entry, &res); err != nil {
			log.Fatalln(err)
		}
		return res.Data
	} else {
		fmt.Println("Fetching....")
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

		var res Response

		if err := json.Unmarshal(body, &res); err != nil {
			defer log.Fatalln(err)
			var body_json map[string]interface{}
			json.Unmarshal(body, &body_json)
			fmt.Println(body_json["code"])
			fmt.Println(body_json["message"])
			if !testing_mode {
				main()
			}

		}

		if res.Code == 200 {
			if !testing_mode {
				err = cache.Set(key, body)
				if err != nil {
					log.Fatalln(err)
				}
			}
			return res.Data
		} else {
			fmt.Println(res.Status)
			return Data{}
		}

	}

}
