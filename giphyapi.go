package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const randomGIFURL = "https://api.giphy.com/v1/gifs/random"

// GIFInfo is the relevant GIF information
type GIFInfo struct {
	Status   string `json:"status"`
	EmbedURL string `json:"embed_url"`
}

// GetRandomGIF fetches random gif data from giphy
func GetRandomGIF(giphyAPIKey string, tag string, rating string) (GIFInfo, error) {
	fmt.Println("Fetching random gif")
	req, err := http.NewRequest("GET", randomGIFURL, nil)
	if err != nil {
		fmt.Println(err)
		return GIFInfo{Status: "API Error"}, errors.New("Error assembling request")
	}

	q := req.URL.Query()
	q.Add("api_key", giphyAPIKey)
	// q.Add("tag", tag)
	q.Add("rating", rating)
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error getting data from api: %v\n", err)
		return GIFInfo{Status: "API Error"}, errors.New("Error getting data from api")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading data from api: %v\n", err)
		return GIFInfo{Status: "API Error"}, errors.New("Error reading data from api")
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Error parsing json")
		return GIFInfo{Status: "API Error"}, errors.New("Error parsing json")
	}

	if data["message"] != nil {
		fmt.Printf("Api returned error message: %v\n", data["message"])
		return GIFInfo{Status: "API Error"}, errors.New("API reported error")
	}

	metaData := data["meta"].(map[string]interface{})
	apiStatus := int(metaData["status"].(float64))
	if apiStatus != 200 {
		fmt.Printf("API returned error")
		return GIFInfo{Status: "API Error"}, errors.New("API Error")
	}

	gifObj := data["data"].(map[string]interface{})
	embedURL := gifObj["embed_url"].(string)

	gif := GIFInfo{EmbedURL: embedURL, Status: "API OK"}

	return gif, nil
}
