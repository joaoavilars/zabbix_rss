package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Config struct to hold the URL from the config file
type Config struct {
	URL string `json:"url"`
}

// RSS struct to unmarshal the XML data
type RSS struct {
	Channel Channel `xml:"channel"`
}

// Channel struct to unmarshal the XML data
type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Items       []Item `xml:"item"`
}

// Item struct to unmarshal the XML data
type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	GUID        string `xml:"guid"`
}

func parseRSSFeed(url string) ([]Item, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rss RSS
	if err := xml.Unmarshal(bytes, &rss); err != nil {
		return nil, err
	}

	return rss.Channel.Items, nil
}

func readConfig(filePath string) (Config, error) {
	var config Config
	file, err := os.Open(filePath)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}

func main() {
	config, err := readConfig("config.json")
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	items, err := parseRSSFeed(config.URL)
	if err != nil {
		fmt.Println("Error parsing RSS feed:", err)
		return
	}

	// Convert items to JSON
	jsonData, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		return
	}

	fmt.Println(string(jsonData))
}
