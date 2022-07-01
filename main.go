package main

import (
	// On importe le packtage fmt

	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// type champions struct {
// 	ID       int    `json:"id"`
// 	Nickname string `json:"nickname"`
// 	Powers   string `json:"powers"`
// 	Tools    string `json:"tools"`
// }

type ItemsResponse struct {
	Items []Brawler `json:"items"`
}

type Brawler struct {
	ID         int     `json:"id"`
	Name       string  `json:"Name"`
	StarPowers []Power `json:"starPowers"`
	Gadgets    []Tool  `json:"gadgets"`
}

type Power struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Tool struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Champion struct {
	Id       int     `json:"id"`
	NickName string  `json:"nickname"`
	Powers   []Power `json:"powers"`
	Tools    []Tool  `json:"tools"`
}

func getChampions(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getChampionsFromBrawlers())
}

func getChampionsFromBrawlers() []Champion {
	// get response from external API through resty library
	client := resty.New()
	var items ItemsResponse
	// var starpowers Brawler
	// var gadgets Brawler

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken("eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiIsImtpZCI6IjI4YTMxOGY3LTAwMDAtYTFlYi03ZmExLTJjNzQzM2M2Y2NhNSJ9.eyJpc3MiOiJzdXBlcmNlbGwiLCJhdWQiOiJzdXBlcmNlbGw6Z2FtZWFwaSIsImp0aSI6IjljYWJlNjNiLWI4ZTQtNGZjNy1iMTBlLWVhMTExOTBmMGJkZSIsImlhdCI6MTY1NTcxMzcyMiwic3ViIjoiZGV2ZWxvcGVyL2M3OTZmM2MxLTE4ZTktNjE0YS0wY2M3LWMwMWIyNTQ2ZDViYiIsInNjb3BlcyI6WyJicmF3bHN0YXJzIl0sImxpbWl0cyI6W3sidGllciI6ImRldmVsb3Blci9zaWx2ZXIiLCJ0eXBlIjoidGhyb3R0bGluZyJ9LHsiY2lkcnMiOlsiMTc4LjIwLjUwLjIwOSJdLCJ0eXBlIjoiY2xpZW50In1dfQ.pbdssmUbaGtoUTKgNOPmspo2gOdtZt7wrgQj3JA1xP8StdKlsg0tAz1iDygAjupUuhjE7ZoHDctNI05vJgwlsw").
		Get("https://api.brawlstars.com/v1/brawlers")

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// unmarshall struct
	json.Unmarshal(resp.Body(), &items)
	fmt.Println(items)

	// //sort items
	return sortBrawler(items.Items)

}

func getChampion(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id = string(params["identifier"])

	client := resty.New()
	var brawler Brawler
	// var starpowers Brawler
	// var gadgets Brawler

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken("eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiIsImtpZCI6IjI4YTMxOGY3LTAwMDAtYTFlYi03ZmExLTJjNzQzM2M2Y2NhNSJ9.eyJpc3MiOiJzdXBlcmNlbGwiLCJhdWQiOiJzdXBlcmNlbGw6Z2FtZWFwaSIsImp0aSI6IjljYWJlNjNiLWI4ZTQtNGZjNy1iMTBlLWVhMTExOTBmMGJkZSIsImlhdCI6MTY1NTcxMzcyMiwic3ViIjoiZGV2ZWxvcGVyL2M3OTZmM2MxLTE4ZTktNjE0YS0wY2M3LWMwMWIyNTQ2ZDViYiIsInNjb3BlcyI6WyJicmF3bHN0YXJzIl0sImxpbWl0cyI6W3sidGllciI6ImRldmVsb3Blci9zaWx2ZXIiLCJ0eXBlIjoidGhyb3R0bGluZyJ9LHsiY2lkcnMiOlsiMTc4LjIwLjUwLjIwOSJdLCJ0eXBlIjoiY2xpZW50In1dfQ.pbdssmUbaGtoUTKgNOPmspo2gOdtZt7wrgQj3JA1xP8StdKlsg0tAz1iDygAjupUuhjE7ZoHDctNI05vJgwlsw").
		Get(fmt.Sprintf("https://api.brawlstars.com/v1/brawlers/%s", id))

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// unmarshall struct
	json.Unmarshal(resp.Body(), &brawler)
	fmt.Println(brawler)

	// //sort items

	newPow := brawler.StarPowers

	newTo := brawler.Gadgets

	newChamp := Champion{
		Id:       brawler.ID,
		NickName: brawler.Name,
		Powers:   newPow,
		Tools:    newTo,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newChamp)

}

func searchChampions(w http.ResponseWriter, r *http.Request) {
	nickname := r.FormValue("nickname")

	champions := getChampionsFromBrawlers()

	result := make([]Champion, 0)

	for i := range champions {
		champion := champions[i]

		if strings.EqualFold(champion.NickName, nickname) {
			result = append(result, champion)
		}

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

}

func sortBrawler(brawlers []Brawler) []Champion {
	var champions []Champion
	// 	for each brawler we create a champion with the same name as nickname and add it
	// 	to the champions lists
	for i := range brawlers {
		newPow := brawlers[i].StarPowers

		newTo := brawlers[i].Gadgets

		// create a champion
		newChamp := Champion{
			Id:       brawlers[i].ID,
			NickName: brawlers[i].Name,
			Powers:   newPow,
			Tools:    newTo,
		}

		// 		add a champion to champions array
		champions = append(champions, newChamp)
	}
	fmt.Println(champions)
	return champions
}

func main() {

	route := mux.NewRouter()
	route.HandleFunc("/champions", getChampions).Methods("GET")
	route.HandleFunc("/champions/{identifier:[0-9]+}", getChampion).Methods("GET")
	route.HandleFunc("/champions/search", searchChampions).Queries("nickname", "{nickname}").Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, handlers.CORS()(route))))
}
