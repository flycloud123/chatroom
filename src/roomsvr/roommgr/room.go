package roommgr

import (
	"fmt"
	"roomsvr/message"
	"roomsvr/stats"
	"sync"
	"time"
)

type Room struct {
	RoomId int

	sessionMutex sync.RWMutex
	Sessions map[*Session]bool
	user2Sessions map[string]*Session



	MaxLatestMsg int
	MsgContainer message.MessageContainer

	RecvBuf  chan *message.TextMessage
}

func (self *Room) addSession(session *Session) {
	self.sessionMutex.Lock()
	self.Sessions[session] = true
	self.user2Sessions[session.User] = session
	self.sessionMutex.Unlock()
}

func (self *Room) removeSession(session *Session) {
	self.sessionMutex.Lock()
	delete(self.Sessions, session)
	delete(self.user2Sessions, session.User)
	self.sessionMutex.Unlock()
}

func (self *Room) GetSessionByUser(user string) *Session {
	self.sessionMutex.RLock()
	ses, _ := self.user2Sessions[user]
	self.sessionMutex.RUnlock()
	return ses
}

func (self *Room) Run()  {
	for {
		select {
		case msg, ok := <- self.RecvBuf:
			if ok {
				self.broadcastMsg(msg)
			}
		}
	}
}

func (self *Room) broadcastMsg(msg *message.TextMessage) {
	self.sessionMutex.RLock()
	for session, _ := range self.Sessions {
		session.PostMsg(msg)
	}
	self.sessionMutex.RUnlock()
}

func (self *Room) Join(session *Session) {
	if session == nil {
		return
	}


	self.addSession(session)
	session.room = self
	session.EnterStamp = time.Now().Unix()

	//加入房间后,推送最近的历史消息
	latestMsgs := self.MsgContainer.FetchLatestMessagesByCount(self.MaxLatestMsg)
	if latestMsgs != nil {
		for _, msg := range latestMsgs {
			session.PostMsg(msg)
		}
	}

	textMessage := message.TextMessage{
		Msg: fmt.Sprintf("%s join room", session.User),
		Sender: "SYSTEM",
		Stamp: time.Now().Unix(),
	}
	self.PostMsg(&textMessage)
}

func (self *Room) Exit(session *Session) {
	if session != nil {
		self.removeSession(session)
		textMessage := message.TextMessage{
			Msg: fmt.Sprintf("%s leave room", session.User),
			Sender: "SYSTEM",
			Stamp: time.Now().Unix(),
		}
		self.PostMsg(&textMessage)
	}
}

func (self *Room) PostMsg(msg *message.TextMessage) {
	self.RecvBuf <- msg
	if msg.Sender != "SYSTEM" {
		self.MsgContainer.StoreMessage(msg)
	}
}

func (self *Room) StatsUser(session *Session, user string) {
	target := self.GetSessionByUser(user)
	if target != nil {
		nt := int(time.Now().Unix() - target.EnterStamp)
		day := nt / 86400
		hour := (nt % 86400) / 3600
		min := (nt % 86400 % 3600) / 60
		sec := nt % 86400 % 3600 % 60
		textMessage := message.TextMessage{
			Msg: fmt.Sprintf("%02dd %02dh %02dm %02ds", day, hour, min, sec),
			Sender: "SYSTEM",
		}
		session.PostMsg(&textMessage)
	} else {
		textMessage := message.TextMessage{
			Msg: fmt.Sprintf("%s not exist", user),
			Sender: "SYSTEM",
		}
		session.PostMsg(&textMessage)
	}

}

func (self *Room) StatsPopular(session *Session, latestStamp int64) {
	if session == nil {
		return
	}

	latestMsgs := self.MsgContainer.FetchLatestMessagesByTime(latestStamp)
	var texts []string
	if latestMsgs != nil {
		for _, msg := range latestMsgs {
			texts = append(texts, msg.Msg)
		}
	}
	popular := stats.FindPopularWord(texts)
	text := ""
	if popular == "" {
		text = "No popular word now"
	} else {
		text = "popular word is " + popular
	}
	textMessage := message.TextMessage{
		Msg: text,
		Sender: "SYSTEM",
	}
	session.PostMsg(&textMessage)
}
