package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/mmcdole/gofeed"
)

func main() {
	fmt.Println("Starting...")
	freq := time.Minute * 2
	for tick := range time.Tick(freq) {
		parseFeed(tick)
	}
}

func parseFeed(start time.Time) {
	url := "https://rpilocator.com/feed.rss"
	fmt.Printf("Fetching feed at %v\n", start.Format("15:04:05"))
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(url)
	for _, item := range feed.Items {
		if item.PublishedParsed.After(start) {
			message := formMessage(item)
			sendMessage(message)
			fmt.Println(message)
		}
	}
}

func formMessage(item *gofeed.Item) string {
	message := fmt.Sprintf("%v\n%v", item.Title, item.Link)
	return message
}

func sendMessage(message string) {
	baseUrl := fmt.Sprintf("https://api.telegram.org/bot%v/sendMessage", os.Getenv("API_TOKEN"))
	v := url.Values{}
	v.Set("chat_id", os.Getenv("CHAT_ID"))
	v.Set("text", message)
	perform := baseUrl + "?" + v.Encode()
	http.Get(perform)
}
