package gss

import (
	"errors"
	"sync"
)

type ChannelManager interface {
	Subscribe(channelName, connID, userID string)
	Unsubscribe(connID, userID string) error
	GetChannelSubscriptions(channelName string) map[string]bool
	GetUserSubscriptions(channelName string, userIDs []string) map[string]bool
}

type channelManager struct {
	ChannelManager
	connIDByUserMap    map[string]*sync.Map
	connIDByChannelMap map[string]*sync.Map
}

func NewChannelManager() *channelManager {
	return &channelManager{
		connIDByUserMap:    map[string]*sync.Map{},
		connIDByChannelMap: map[string]*sync.Map{},
	}
}

func (m *channelManager) GetMapsByUser(userID string) (*sync.Map, error) {
	idmap, exists := m.connIDByUserMap[userID]
	if !exists {
		return nil, errors.New("userID: " + userID + " not registered")
	}
	return idmap, nil
}

func (m *channelManager) GetMapsByChannel(channelName string) (*sync.Map, error) {
	idmap, exists := m.connIDByChannelMap[channelName]
	if !exists {
		return nil, errors.New("channelName: " + channelName + " not registered")
	}
	return idmap, nil
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

func (m *channelManager) Unsubscribe(connId, userId string) error {
	userStore, err := m.GetMapsByUser(userId)
	if err != nil {
		return err
	}
	userStore.Delete(connId)
	if !keyExists(userStore) {
		delete(m.connIDByUserMap, userId)
	}
	for chname, store := range m.connIDByChannelMap {
		store.Delete(connId)
		if !keyExists(store) {
			delete(m.connIDByChannelMap, chname)
		}
	}

	return nil
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
