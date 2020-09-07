package roommgr

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"roomsvr/message"
	"roomsvr/wordfilter"
	"strconv"
	"strings"
	"time"
)

type Session struct {
	User string

	Conn *websocket.Conn
	SendBuf chan *message.TextMessage
	stopCh chan int
	EnterStamp int64

	room *Room
}

func (self *Session) PostMsg(msg *message.TextMessage) {
	self.SendBuf <- msg
}

func (self *Session) Run() {
	defer func() {
		self.Conn.Close()
		if self.room != nil {
			self.room.Exit(self)
		}
		RoomMgr.Exit(self)
	}()

	self.stopCh = make(chan int)
	go self.sendMsg()

	for {
		_, raw, err := self.Conn.ReadMessage()
		if err != nil {
			close(self.stopCh)
			return
		}

		self.handler(raw)
	}
}

func (self *Session) handler(rawmsg []byte) {
	datas := make(map[string]string)
	err := json.Unmarshal(rawmsg, &datas)
	if err == nil {
		if cmd, ok := datas["cmd"]; ok {
			switch cmd {
			case "create_room": self.createRoomHandler(datas)
			case "join_room": self.joinRoomHandler(datas)
			case "send_msg": self.sendMsgHandler(datas)
			}
		}
	}
}

func (self *Session) createRoomHandler(datas map[string]string) {
	if user, ok := datas["user"]; ok {
		self.User = user
		RoomMgr.CreateRoom(self)
	}
}

func (self *Session) joinRoomHandler(datas map[string]string) {
	if user, ok := datas["user"]; ok {
		self.User = user

		if roomId, ok := datas["roomid"]; ok {
			roomIdInt, err := strconv.Atoi(roomId)
			if err == nil {
				RoomMgr.JoinRoom(roomIdInt, self)
			}
		}
	}
}

func (self *Session) forwardProcess(text string) bool {
	if strings.Index(text, "/stats") >= 0 {
		datas := strings.Split(text, " ")
		if len(datas) == 2 {
			self.room.StatsUser(self, datas[1])
			return true
		}
	} else if strings.Index(text, "/popular") >= 0 {
		datas := strings.Split(text, " ")
		if len(datas) == 2 {
			second, err := strconv.Atoi(datas[1])
			if err == nil {
				latestStamp := time.Now().Unix() - int64(second)
				self.room.StatsPopular(self, latestStamp)
				return true
			}

		}
	}
	return false
}

func (self *Session) sendMsgHandler(datas map[string]string) {
	if self.room != nil {
		if text, ok := datas["text"]; ok {
			//关键词过滤
			raw := wordfilter.ReplaceDirty([]byte(text))
			text = string(raw)
			if self.forwardProcess(text) {
				return
			}

			textMsg := message.TextMessage{
				Msg: text,
				Sender: self.User,
				Stamp: time.Now().Unix(),
			}
			self.room.PostMsg(&textMsg)
		}
	}
}

func (self *Session) sendMsg() {
	for {
		select {
		case <- self.stopCh:
			return
		case msg, ok := <- self.SendBuf:
			if ok {
				raw := []byte(msg.Sender + ": ")	//带上用户名,方便查看
				raw = append(raw, msg.Msg...)
				self.Conn.WriteMessage(websocket.TextMessage, raw)
			}
		}
	}
}

func NewSession(conn *websocket.Conn, sendBufSize int, timeout int) *Session {
	session := &Session{
		Conn: conn,
		SendBuf: make(chan *message.TextMessage, sendBufSize),
	}
	return session
}