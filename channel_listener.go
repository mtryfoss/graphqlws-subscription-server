package graphqlws_subscription_server

import (
	"sync"
)

type Listener struct {
	channelMapMutex    *sync.RWMutex
	userMapMutex       *sync.RWMutex
	connIdByUserMap    map[string]map[string]bool
	connIdByChannelMap map[string]map[string]bool
	dummyLabel         string
}

func NewListener(dummy string) *Listener {
	return &Listener{
		channelMapMutex:    &sync.RWMutex{},
		userMapMutex:       &sync.RWMutex{},
		connIdByUserMap:    map[string]map[string]bool{},
		connIdByChannelMap: map[string]map[string]bool{},
		dummyLabel:         dummy,
	}
}

type ConnectionsByChannel map[string]bool

func (l *Listener) Subscribe(channel, connId, userId string) {
	l.channelMapMutex.RLock()
	if connList, exists := l.connIdByChannelMap[channel]; exists {
		connList[connId] = true
		l.connIdByChannelMap[channel] = connList
	} else {
		l.connIdByChannelMap[channel] = map[string]bool{connId: true}
	}
	l.channelMapMutex.RUnlock()
	if userId == l.dummyLabel {
		return
	}
	l.userMapMutex.RLock()
	if connList, exists := l.connIdByUserMap[userId]; exists {
		connList[connId] = true
		l.connIdByUserMap[userId] = connList
	} else {
		l.connIdByUserMap[userId] = map[string]bool{connId: true}
	}
	l.userMapMutex.RUnlock()
}

func (l *Listener) Unsubscribe(connId, userId string) {
	l.channelMapMutex.Lock()
	connIds := []string{connId}
	if userId != l.dummyLabel {
		l.userMapMutex.Lock()
		for cid, _ := range l.connIdByUserMap[userId] {
			if cid != connId {
				connIds = append(connIds, cid)
			}
		}
		delete(l.connIdByUserMap, userId)
		l.userMapMutex.Unlock()
	}
	for channel, connList := range l.connIdByChannelMap {
		for _, cid := range connIds {
			delete(connList, cid)
		}
		l.connIdByChannelMap[channel] = connList
	}
	l.channelMapMutex.Unlock()
}

func (l *Listener) GetChannelSubscribers(channel string) []string {
	listenerConns := []string{}
	l.channelMapMutex.RLock()
	for cid, _ := range l.connIdByChannelMap[channel] {
		listenerConns = append(listenerConns, cid)
	}
	l.channelMapMutex.RUnlock()
	return listenerConns
}

func (l *Listener) GetUserSubscribers(userIds []string) []string {
	listenerConns := []string{}
	l.userMapMutex.RLock()
	for _, uid := range userIds {
		for cid, _ := range l.connIdByUserMap[uid] {
			listenerConns = append(listenerConns, cid)
		}
	}
	l.userMapMutex.RUnlock()
	return listenerConns
}
