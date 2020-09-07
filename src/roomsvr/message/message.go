package message

import (
	"sync"
)

type TextMessage struct {
	Msg    string
	Sender string
	Stamp  int64
}

type MessageContainer interface {
	StoreMessage(textMessage *TextMessage)
	FetchLatestMessagesByCount(count int) []*TextMessage
	FetchLatestMessagesByTime(latestStamp int64) []*TextMessage
}

//这里只实现了一个简单的消息存储在内存中,如果要求持久化,可以保存在数据库中
type MemMsgContainer struct {
	LatestMsgs []*TextMessage
	LatestMutex sync.RWMutex
}

func (self *MemMsgContainer) StoreMessage(msg *TextMessage) {
	self.LatestMutex.Lock()

	self.LatestMsgs = append(self.LatestMsgs, msg)

	self.LatestMutex.Unlock()
}

func (self *MemMsgContainer) FetchLatestMessagesByCount(count int) []*TextMessage {
	messages := make([]*TextMessage, 0, count)

	self.LatestMutex.RLock()
	if self.LatestMsgs != nil {
		idx := 0
		if count < len(self.LatestMsgs) {
			idx = len(self.LatestMsgs) - count
		}
		for ; idx < len(self.LatestMsgs); idx++ {
			messages = append(messages, self.LatestMsgs[idx])
		}
	}
	self.LatestMutex.RUnlock()

	return messages
}

func (self *MemMsgContainer) FetchLatestMessagesByTime(latestStamp int64) []*TextMessage {
	var messages []*TextMessage

	self.LatestMutex.RLock()
	if self.LatestMsgs != nil && len(self.LatestMsgs) > 0 {
		for idx := len(self.LatestMsgs) - 1; idx >= 0; idx-- {
			if self.LatestMsgs[idx].Stamp >= latestStamp {
				messages = append(messages, self.LatestMsgs[idx])
			} else {
				break
			}
		}
	}
	self.LatestMutex.RUnlock()

	for idx := 0; idx < len(messages) / 2; idx++ {
		msg := messages[idx]
		messages[idx] = messages[len(messages) - 1 - idx]
		messages[len(messages) - 1 - idx] = msg
	}
	return messages
}