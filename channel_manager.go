package gss

import (
	"sync"
)

type ChannelManager interface {
	Subscribe(string, string, string)
	Unsubscribe(string, string)
	GetChannelSubscriptions(string) map[string]bool
	GetUserSubscriptions(string, []string) map[string]bool
}

type channelManager struct {
	ChannelManager
	connIDByUserMap    map[string]*sync.Map
	connIDByChannelMap map[string]*sync.Map
}

func NewChannelManager() ChannelManager {
	return &channelManager{
		connIDByUserMap:    map[string]*sync.Map{},
		connIDByChannelMap: map[string]*sync.Map{},
	}
}

func (m *channelManager) Subscribe(channel, connId, userId string) {
	if connList, exists := m.connIDByChannelMap[channel]; exists {
		connList.Store(connId, true)
	} else {
		store := &sync.Map{}
		store.Store(connId, true)
		m.connIDByChannelMap[channel] = store
	}
	if connList, exists := m.connIDByUserMap[userId]; exists {
		connList.Store(connId, true)
	} else {
		store := &sync.Map{}
		store.Store(connId, true)
		m.connIDByUserMap[userId] = store
	}
}

func keyExists(m *sync.Map) bool {
	cnt := 0
	m.Range(func(k, v interface{}) bool {
		cnt++
		return false
	})
	return cnt > 0
}

func (m *channelManager) Unsubscribe(connId, userId string) {
	connIds := []string{connId}
	if store, exists := m.connIDByUserMap[userId]; exists {
		store.Range(func(k, v interface{}) bool {
			connIds = append(connIds, k.(string))
			return true
		})
		delete(m.connIDByUserMap, userId)
	}
	for chname, store := range m.connIDByChannelMap {
		for _, cid := range connIds {
			store.Delete(cid)
		}
		if !keyExists(store) {
			delete(m.connIDByChannelMap, chname)
		}
	}
}

func (m *channelManager) GetChannelSubscriptions(channel string) map[string]bool {
	connIds := map[string]bool{}
	if connList, exists := m.connIDByChannelMap[channel]; exists {
		connList.Range(func(k, v interface{}) bool {
			connIds[k.(string)] = true
			return true
		})
	}
	return connIds
}

func (m *channelManager) GetUserSubscriptions(channel string, userIds []string) map[string]bool {
	connIds := map[string]bool{}
	if connList, exists := m.connIDByChannelMap[channel]; exists {
		connList.Range(func(k, v interface{}) bool {
			connIds[k.(string)] = true
			return true
		})
	}
	userConnIds := map[string]bool{}
	for _, uid := range userIds {
		if connList, exists := m.connIDByUserMap[uid]; exists {
			connList.Range(func(k, v interface{}) bool {
				connID := k.(string)
				if _, exists := connIds[connID]; exists {
					userConnIds[connID] = true
				}
				return true
			})
		}
	}
	return userConnIds
}
