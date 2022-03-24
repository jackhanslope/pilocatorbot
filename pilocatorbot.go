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
	scheduled_start := time.Now()
	freq := time.Minute * 2
	for {
		parseFeed(scheduled_start)
		diff := time.Now().Sub(scheduled_start)
		time.Sleep(freq - diff)
		scheduled_start = scheduled_start.Add(freq)
	}
}

func parseFeed(scheduled_start time.Time) {
	url := "https://rpilocator.com/feed.rss"
	fmt.Printf("Fetching feed at %v\n", scheduled_start.Format("15:04:05"))
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(url)
	for _, item := range feed.Items {
		if item.PublishedParsed.After(scheduled_start) {
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
