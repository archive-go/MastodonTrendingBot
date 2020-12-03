package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/MakeGolangGreat/mastodon-go"
	"github.com/ansel1/merry"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var token string
var domains []string

func loadConfig() {
	err := godotenv.Load()
	merry.Wrap(err)

	token = os.Getenv("TOKEN")
	fmt.Println(token)
	domains = strings.Fields(os.Getenv("DOMAIN"))
	fmt.Println(domains)
}

func init() {
	loadConfig()
}

func main() {
	for _, instance := range domains {
		listen(instance)
	}
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

			process(status)
		}
	}
}

func process(status mastodon.Status) {
	fmt.Println(status.Content)
}
