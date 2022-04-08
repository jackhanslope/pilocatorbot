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
		if err = parseFeed(tick, conf); err != nil {
			log.Println(err)
		}
	}
}

func loadConfig() (conf config, err error) {
	err = env.Parse(&conf)
	return
}

func parseFeed(start time.Time, conf config) (err error) {
	url := conf.RssUrl
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return
	}
	for _, item := range feed.Items {
		if item.PublishedParsed.After(start.Add(-conf.UpdateFreq)) {
			message := formMessage(item)
			sendMessage(message, conf)
			log.Println(message)
		}
	}
	return
}

func formMessage(item *gofeed.Item) (message string) {
	url, _ := url.Parse(item.Link)
	message = fmt.Sprintf("%v %v://%v", item.Title, url.Scheme, url.Host)
	return
}

func sendMessage(message string, conf config) {
	baseUrl := fmt.Sprintf("https://api.telegram.org/bot%v/sendMessage", conf.Token)
	v := url.Values{}
	v.Set("chat_id", conf.ChatId)
	v.Set("text", message)
	perform := baseUrl + "?" + v.Encode()
	http.Get(perform)
}
