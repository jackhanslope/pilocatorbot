package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/mmcdole/gofeed"
)

type config struct {
	Token      string        `env:"PILOCBOT_TOKEN,notEmpty"`
	ChatId     string        `env:"PILOCBOT_CHAT_ID,notEmpty"`
	RssUrl     string        `env:"PILOCBOT_RSS_URL,notEmpty"`
	UpdateFreq time.Duration `env:"PILOCBOT_UPDATE_FREQ" envDefault:"2m"`
}

func main() {
	conf, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Fetching at %v every %v\n", conf.RssUrl, conf.UpdateFreq)
	for tick := range time.Tick(conf.UpdateFreq) {
		parseFeed(tick, conf)
	}
}

func loadConfig() (conf config, err error) {
	err = env.Parse(&conf)
	if err != nil {
		return
	}
	return
}

func parseFeed(start time.Time, conf config) {
	url := conf.RssUrl
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(url)
	for _, item := range feed.Items {
		if item.PublishedParsed.After(start.Add(-conf.UpdateFreq)) {
			message := formMessage(item)
			sendMessage(message, conf)
			log.Println(message)
		}
	}
}

func formMessage(item *gofeed.Item) string {
	message := fmt.Sprintf("%v\n%v", item.Title, item.Link)
	return message
}

func sendMessage(message string, conf config) {
	baseUrl := fmt.Sprintf("https://api.telegram.org/bot%v/sendMessage", conf.Token)
	v := url.Values{}
	v.Set("chat_id", conf.ChatId)
	v.Set("text", message)
	perform := baseUrl + "?" + v.Encode()
	http.Get(perform)
}
