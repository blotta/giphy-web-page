package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var config Config

// Config represents information from a config.json file
type Config struct {
	GiphyAPIKey string `json:"giphy_apikey"`
}

// PageData is the data for the html template
type PageData struct {
	Greeting string
	Gif      GIFInfo
}

func main() {

	if err := readConfig(&config); err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	// Static files
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/")))).Methods("GET")

	// Routes //
	////////////

	// Home
	r.HandleFunc("/", Index)

	// randomgif api for ajax
	r.HandleFunc("/randomgif", RandomGIF)

	http.ListenAndServe(":8080", r)
}

func readConfig(config *Config) error {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	if err := json.Unmarshal(byteValue, config); err != nil {
		return err
	}
	return nil
}

// Index handles request for root
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request from :", r.RemoteAddr, " to", r.RequestURI)
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	gif, err := GetRandomGIF(config.GiphyAPIKey, "", "G")
	if err != nil {
		fmt.Printf("Error getting gif: %v", err)
	}
	data := PageData{
		Greeting: "Hello, there!",
		Gif:      gif,
	}

	tmpl.Execute(w, data)
}

// RandomGIF returns a json with a random gif information
func RandomGIF(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	gif, _ := GetRandomGIF(config.GiphyAPIKey, "", "G")

	respJSON, err := json.Marshal(gif)
	if err != nil {
		fmt.Println("Error converting to json")
		return
	}

	fmt.Fprint(w, string(respJSON))
}
