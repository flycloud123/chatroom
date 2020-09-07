package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var user = flag.String("user", "Roland", "user name")
var cmd = flag.String("cmd", "create_room/join_room", "join room")
var roomId = flag.String("roomid", "1", "room id")

func main() {
	flag.Parse()
	log.SetFlags(0)

	fmt.Printf("addr: %s, user: %s, cmd: %s, roomId: %s\n", *addr, *user, *cmd, *roomId)

	inputMsgs := make(chan string, 100)
	go func() {
		inputReader := bufio.NewReader(os.Stdin)
		for {
			text, err := inputReader.ReadString('\n')
			if err == nil {
				text = strings.Replace(text, "\n", "", 1)
				inputMsgs <- text
			}
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}

	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	msg := make(map[string]string)
	msg["user"] = *user
	msg["cmd"] = *cmd
	if msg["cmd"] == "join_room" {
		msg["roomid"] = *roomId
	}
	raw, _ := json.Marshal(msg)
	c.WriteMessage(websocket.TextMessage, raw)

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("recv:", err)
				return
			}
			log.Printf("%s", message)
		}
	}()

	for {
		select {
		case <-done:
			return
		case text, ok := <- inputMsgs:
			if ok {
				msg["cmd"] = "send_msg"
				msg["text"] = text
				raw, _ := json.Marshal(msg)
				err := c.WriteMessage(websocket.TextMessage, raw)
				if err != nil {
					log.Println("write:", err)
					return
				}
			}
		case <-interrupt:
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

