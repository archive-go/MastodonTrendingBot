package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/MakeGolangGreat/mastodon-go"
	"github.com/PuerkitoBio/goquery"
	"github.com/ansel1/merry"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
)

var token string
var domains []string

func init() {
	loadConfig()
	if len(domains) == 0 || token == "" {
		color.Red("Missing domains/token")
		log.Fatal()
	}
}

func main() {
	defer db.Close()

	for _, instance := range domains {
		cronJob(instance)
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

			fmt.Println("Status ID:", status.ID, status.Account.UserName)
			process(status, domain)
		}
	}
}

func process(status mastodon.Status, domain string) {
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

		count := get(s.Text())
		set(s.Text(), count+1, domain)
	})

}

func loadConfig() {
	err := godotenv.Load()
	merry.Wrap(err)

	token = os.Getenv("TOKEN")
	domains = strings.Fields(os.Getenv("DOMAINS"))
}

func publish(domain string) {
	tags := getAll()
	toot := &mastodon.Mastodon{
		Token:  token,
		Domain: "https://" + domain,
	}

	var content string
	for i := 0; i < min(len(tags), 10); i++ {
		content += fmt.Sprintf("%d. %s\n", i+1, tags[i].name)
	}
	var status string
	if len(tags) == 0 {
		status = "今天居然没有一个人使用标签，我一个Bot能怎么办？\n\n标签用起来呀，朋友们～"
	} else {
		status = fmt.Sprintf("%s 的兄弟姐妹们晚上好，我来为大家播报本站今日热门（排名有先后）：\n\n%s\n\n为了让本Bot能如实还原热门趋势情况，还请大家多使用“标签”功能！", domain, content)
	}

	_, err := toot.SendStatuses(&mastodon.StatusParams{
		Status:     status,
		MediaIds:   "[]",
		Poll:       "[]",
		Visibility: "private",
		Sensitive:  false,
	})
	merry.Wrap(err)
}

func cronJob(domain string) {
	// 每天的零点零分零秒
	spec := "0 0 0,12 * * *"
	secondsEastOfUTC := int((8 * time.Hour).Seconds())
	beijing := time.FixedZone("Beijing Time", secondsEastOfUTC)
	c := cron.NewWithLocation(beijing)
	// c := cron.New()
	c.AddFunc(spec, func() {
		fmt.Print(3)
		publish(domain)
	})
	c.Start()
}
