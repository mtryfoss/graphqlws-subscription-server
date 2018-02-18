package graphqlws_subscription_server

import (
	"sync"
)

type ConnectionsByID map[string]bool

type Listener struct {
	channelMapMutex    *sync.RWMutex
	userMapMutex       *sync.RWMutex
	connIDByUserMap    map[string]ConnectionsByID
	connIDByChannelMap map[string]ConnectionsByID
	dummyLabel         string
}

func NewListener(dummy string) *Listener {
	return &Listener{
		channelMapMutex:    &sync.RWMutex{},
		userMapMutex:       &sync.RWMutex{},
		connIDByUserMap:    map[string]ConnectionsByID{},
		connIDByChannelMap: map[string]ConnectionsByID{},
		dummyLabel:         dummy,
	}
}

func (l *Listener) Subscribe(channel, connId, userId string) {
	l.channelMapMutex.RLock()
	if connList, exists := l.connIDByChannelMap[channel]; exists {
		connList[connId] = true
		l.connIDByChannelMap[channel] = connList
	} else {
		l.connIDByChannelMap[channel] = ConnectionsByID{connId: true}
	}
	l.channelMapMutex.RUnlock()
	if userId == l.dummyLabel {
		return
	}
	l.userMapMutex.RLock()
	if connList, exists := l.connIDByUserMap[userId]; exists {
		connList[connId] = true
		l.connIDByUserMap[userId] = connList
	} else {
		l.connIDByUserMap[userId] = ConnectionsByID{connId: true}
	}
	l.userMapMutex.RUnlock()
}

func (l *Listener) Unsubscribe(connId, userId string) {
	l.channelMapMutex.Lock()
	connIds := []string{connId}
	if userId != l.dummyLabel {
		l.userMapMutex.Lock()
		for cid, _ := range l.connIDByUserMap[userId] {
			if cid != connId {
				connIds = append(connIds, cid)
			}
		}
		delete(l.connIDByUserMap, userId)
		l.userMapMutex.Unlock()
	}
	for channel, connList := range l.connIDByChannelMap {
		for _, cid := range connIds {
			delete(connList, cid)
		}
		l.connIDByChannelMap[channel] = connList
	}
	l.channelMapMutex.Unlock()
}

func (l *Listener) GetChannelSubscribers(channel string) []string {
	listenerConns := []string{}
	l.channelMapMutex.RLock()
	for cid, _ := range l.connIDByChannelMap[channel] {
		listenerConns = append(listenerConns, cid)
	}
	l.channelMapMutex.RUnlock()
	return listenerConns
}

func (l *Listener) GetUserSubscribers(userIds []string) []string {
	listenerConns := []string{}
	l.userMapMutex.RLock()
	for _, uid := range userIds {
		for cid, _ := range l.connIDByUserMap[uid] {
			listenerConns = append(listenerConns, cid)
		}
	}
	l.userMapMutex.RUnlock()
	return listenerConns
}
