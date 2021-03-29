package main

import (
	"encoding/json"
	"log"
	"net/http"
)

//struct para armazenar informações retornadas pelo catApi
type CatResponse []struct {
	Url string `json:"url"`
}

//conjunto de structs para armazenar informações retornadas pela api do reddit
type redditResponse struct {
	Kind string      `json:"kind"`
	Data dataReponse `json:"data"`
}
type dataReponse struct {
	ModHash  string           `json:"modhash"`
	Dist     int              `json:"dist"`
	Children []childrenReddit `json:"children"`
	After    string           `json:"after"`
	Before   string           `json:"before"`
}
type childrenReddit struct {
	Kind string             `json:"kind"`
	Data dataChildrenReddit `json:"data"`
}
type dataChildrenReddit struct {
	Title string `json:"title"`
	Url   string `json:"url"`
	Type  string `json:"post_hint"`
}

func get_catApi_pic_Url() string {
	url := "https://api.thecatapi.com/v1/images/search"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var urlcat CatResponse

	if err := json.NewDecoder(resp.Body).Decode(&urlcat); err != nil {
		//log.Fatal("ooopsss! an error occurred, please try again")
		log.Fatal(err)
	}
	return urlcat[0].Url

}
func get_raww_top3_pics_url() (string, string, string) {
	url := "https://www.reddit.com/r/aww/hot/.json?limit=25"
	cliente := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("User-Agent", "tbot.ld")

	resp, err := cliente.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var rr redditResponse

	if err := json.NewDecoder(resp.Body).Decode(&rr); err != nil {
		//log.Fatal("ooopsss! an error occurred, please try again")
		log.Fatal(err)
	}
	var str1, str2, str3 string = "", "", ""
	for i, j := 0, 0; i < len(rr.Data.Children) && j < 3; i++ {
		if rr.Data.Children[i].Data.Type == "image" {
			switch j {
			case 0:
				j++
				str1 = rr.Data.Children[i].Data.Url
			case 1:
				j++
				str2 = rr.Data.Children[i].Data.Url
			case 2:
				j++
				str3 = rr.Data.Children[i].Data.Url
			}
		}
	}
	return str1, str2, str3
}
