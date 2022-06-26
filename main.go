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
var total int = 0
var limit int = 10

type Response struct {
	Code int  `json:"code"`
	Data Data `json:"data"`
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

	results := get_marvel_data()
	print_results(results)
	nextAction()

}

func print_results(results Data) {
	total = results.Total
	fmt.Println("Total: ", total)
	for _, v := range results.Results {
		thumb := v.Thumbnail
		fmt.Println("Name: ", v.Name)
		fmt.Println("Description: ", v.Description)
		fmt.Println("Image: ", thumb.Path)
		fmt.Println("-------------------------------------")
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
}

func get_marvel_data() Data {
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
		panic(err)
	}

	if res.Code == 200 {
		return res.Data
	} else {
		return Data{}
	}
}
