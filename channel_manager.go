package gss

import (
	"sync"
)

type ChannelManager interface {
	Subscribe(channelName, connID, userID string)
	Unsubscribe(connID, userID string)
	GetChannelSubscriptions(channelName string) map[string]bool
	GetUserSubscriptions(channelName string, userIDs []string) map[string]bool
}

type channelManager struct {
	ChannelManager
	ConnIDByUserMap    map[string]*sync.Map
	ConnIDByChannelMap map[string]*sync.Map
}

func NewChannelManager() ChannelManager {
	return &channelManager{
		ConnIDByUserMap:    map[string]*sync.Map{},
		ConnIDByChannelMap: map[string]*sync.Map{},
	}
}

func (m *channelManager) Subscribe(channel, connId, userId string) {
	if connList, exists := m.ConnIDByChannelMap[channel]; exists {
		connList.Store(connId, true)
	} else {
		store := &sync.Map{}
		store.Store(connId, true)
		m.ConnIDByChannelMap[channel] = store
	}
	if connList, exists := m.ConnIDByUserMap[userId]; exists {
		connList.Store(connId, true)
	} else {
		store := &sync.Map{}
		store.Store(connId, true)
		m.ConnIDByUserMap[userId] = store
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
	if store, exists := m.ConnIDByUserMap[userId]; exists {
		store.Range(func(k, v interface{}) bool {
			connIds = append(connIds, k.(string))
			return true
		})
		delete(m.ConnIDByUserMap, userId)
	}
	for chname, store := range m.ConnIDByChannelMap {
		for _, cid := range connIds {
			store.Delete(cid)
		}
		if !keyExists(store) {
			delete(m.ConnIDByChannelMap, chname)
		}
	}
}

func (m *channelManager) GetChannelSubscriptions(channel string) map[string]bool {
	connIds := map[string]bool{}
	if connList, exists := m.ConnIDByChannelMap[channel]; exists {
		connList.Range(func(k, v interface{}) bool {
			connIds[k.(string)] = true
			return true
		})
	}
	return connIds
}

func (m *channelManager) GetUserSubscriptions(channel string, userIds []string) map[string]bool {
	connIds := map[string]bool{}
	if connList, exists := m.ConnIDByChannelMap[channel]; exists {
		connList.Range(func(k, v interface{}) bool {
			connIds[k.(string)] = true
			return true
		})
	}
	userConnIds := map[string]bool{}
	for _, uid := range userIds {
		if connList, exists := m.ConnIDByUserMap[uid]; exists {
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
