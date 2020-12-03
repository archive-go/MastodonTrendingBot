package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/MakeGolangGreat/mastodon-go"
	"github.com/PuerkitoBio/goquery"
	"github.com/ansel1/merry"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

var token string
var domains []string

func init() {
	loadConfig()
}

func main() {
	defer db.Close()

	cronJob()

	for _, instance := range domains {
		go listen(instance)
	}

	// for {} will use 100% cpu.
	select {}
}

func listen(domain string) {
	addr := "wss://" + domain + "/api/v1/streaming/?stream=public:local"
	ws, _, err := websocket.DefaultDialer.Dial(addr, nil)
	defer ws.Close()
	if err != nil {
		merry.Wrap(err)
	}

	fmt.Printf("成功连接WS，持续监听实例：%s\n", domain)
	for {
		_, body, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		var message stream
		if err := json.Unmarshal(body, &message); err != nil {
			color.Red("解析字符串出错！", err)
		}

		switch message.Event {
		case "update":
			var status mastodon.Status
			err := json.Unmarshal([]byte(message.Payload.(string)), &status)
			if err != nil {
				fmt.Println("err", err.Error())
			}

			process(status, domain)
		}
	}
}

func process(status mastodon.Status, domain string) {
	fmt.Println(domain, status.Content)

	count := get(domain)
	newCount, err := strconv.Atoi(count)
	fmt.Println("domain count is: ", newCount)

	merry.Wrap(err)
	set(domain, newCount+1)
	fmt.Println("setted ", get(domain))

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(status.Content))
	merry.Wrap(err)

	// find tag link
	doc.Find("a.hashtag").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		url, err := url.Parse(href)
		merry.Wrap(err)
		if url.Host != domain {
			// 如果链接是当前实例的，那么不会将其备份
			return
		}

		fmt.Println(domain, s.Text())
	})

}

func loadConfig() {
	err := godotenv.Load()
	merry.Wrap(err)

	token = os.Getenv("TOKEN")
	fmt.Println(token)
	domains = strings.Fields(os.Getenv("DOMAINS"))
	fmt.Println(domains)
}

func publish() {
	fmt.Println("beijing 0 dian")
}

func cronJob() {
	getAll()
	spec := "* 0 * * *"
	secondsEastOfUTC := int((8 * time.Hour).Seconds())
	beijing := time.FixedZone("Beijing Time", secondsEastOfUTC)
	c := cron.New(cron.WithLocation(beijing))
	c.AddFunc(spec, publish)
	c.Start()
}
