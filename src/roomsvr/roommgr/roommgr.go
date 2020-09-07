package roommgr

import (
	"fmt"
	"github.com/gorilla/websocket"
	"roomsvr/message"
	"sync"
	"time"
)

type RoomManager struct {
	nextRoomId int
	rooms      map[int]*Room
	roomMutex  sync.RWMutex

	allSessions map[*Session]bool
	sessionMutex sync.Mutex
}

var RoomMgr *RoomManager

func init() {
	RoomMgr = &RoomManager{
		nextRoomId: 0,
		rooms: make(map[int]*Room),
		allSessions: make(map[*Session]bool),
	}
}

func (self *RoomManager) CreateRoom(session *Session) *Room {
	self.roomMutex.Lock()
	defer self.roomMutex.Unlock()

	self.nextRoomId++
	room := &Room{
		RoomId: self.nextRoomId,
		Sessions: make(map[*Session]bool),
		user2Sessions: make(map[string]*Session),
		MaxLatestMsg: 50,
		RecvBuf: make(chan *message.TextMessage, 1000),
		MsgContainer: &message.MemMsgContainer{},
	}
	self.rooms[room.RoomId] = room

	textMessage := message.TextMessage{
		Msg: fmt.Sprintf("%s create room %d success", session.User, room.RoomId),
		Sender: "SYSTEM",
		Stamp: time.Now().Unix(),
	}
	room.PostMsg(&textMessage)

	room.Join(session)


	go room.Run()
	return room
}

func (self *RoomManager) JoinRoom(roomId int, session *Session) *Room {
	if session == nil {
		return nil
	}

	self.roomMutex.RLock()
	defer self.roomMutex.RUnlock()

	if room, ok := self.rooms[roomId]; ok {
		room.Join(session)
		return room
	}
	return nil
}

func (self *RoomManager) Enter(conn *websocket.Conn) {
	session := NewSession(conn,100, 6)

	self.sessionMutex.Lock()
	self.allSessions[session] = true
	defer self.sessionMutex.Unlock()

	go session.Run()
}

func (self *RoomManager) Exit(session *Session) {
	self.sessionMutex.Lock()
	defer self.sessionMutex.Unlock()

	delete(self.allSessions, session)
}