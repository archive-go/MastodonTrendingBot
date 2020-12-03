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
var domain []string

func loadConfig() {
	err := godotenv.Load()
	merry.Wrap(err)

	token = os.Getenv("TOKEN")
	fmt.Println(token)
	domain = strings.Fields(os.Getenv("DOMAIN"))
	fmt.Println(domain)
}

func init() {
	loadConfig()
}

func main() {
	addr := "wss://" + "bgme.me" + "/api/v1/streaming/?stream=public:local"
	ws, res, err := websocket.DefaultDialer.Dial(addr, nil)
	defer ws.Close()
	if err != nil {
		fmt.Println(res)
		merry.Wrap(err)
	}

	fmt.Println("连接上WS，持续监听")
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

		}
	}
}
