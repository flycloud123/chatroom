package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"roomsvr/roommgr"
	"roomsvr/wordfilter"
	"strings"
)

type WsServer struct {
	listener net.Listener
	addr     string
	upgrade  *websocket.Upgrader
}

var room *roommgr.Room
var prefix string = "/ws"

func NewWsServer() *WsServer {
	ws := new(WsServer)
	ws.addr = "0.0.0.0:8080"
	ws.upgrade = &websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			if r.Method != "GET" {
				fmt.Println("method is not GET")
				return false
			}
			if strings.Index(r.URL.Path, prefix) != 0 {
				fmt.Println("path error")
				return false
			}
			return true
		},
	}
	return ws
}

func (self *WsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.Index(r.URL.Path, prefix) != 0 {
		httpCode := http.StatusInternalServerError
		reasePhrase := http.StatusText(httpCode)
		http.Error(w, reasePhrase, httpCode)
		return
	}

	conn, err := self.upgrade.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("websocket error:", err)
		return
	}

	roommgr.RoomMgr.Enter(conn)
}

func (w *WsServer) Start() (err error) {
	w.listener, err = net.Listen("tcp", w.addr)
	if err != nil {
		fmt.Println("net listen error:", err)
		return
	}

	fmt.Printf("server listen on %s sucess\n", w.addr)

	err = http.Serve(w.listener, w)
	if err != nil {
		fmt.Println("http serve error:", err)
		return
	}
	return nil
}

func main() {
	wordfilter.LoadDirtyFromFile("dirtyword.txt")

	ws := NewWsServer()
	ws.Start()
}
