package gss

import (
	"errors"
	"sync"
)

type SubscribeFilter interface {
	Subscribe(channelName, connID, userID string)
	Unsubscribe(connID, userID string) error
	GetChannelSubscriptions(channelName string) map[string]bool
	GetUserSubscriptions(channelName string, userIDs []string) map[string]bool
}

type subscribeFilter struct {
	SubscribeFilter
	connIDByUserMap    map[string]*sync.Map
	connIDByChannelMap map[string]*sync.Map
}

func NewSubscribeFilter() *subscribeFilter {
	return &subscribeFilter{
		connIDByUserMap:    map[string]*sync.Map{},
		connIDByChannelMap: map[string]*sync.Map{},
	}
}

func (f *subscribeFilter) GetMapsByUser(userID string) (*sync.Map, error) {
	idmap, exists := f.connIDByUserMap[userID]
	if !exists {
		return nil, errors.New("userID: " + userID + " not registered")
	}
	return idmap, nil
}

func (f *subscribeFilter) GetMapsByChannel(channelName string) (*sync.Map, error) {
	idmap, exists := f.connIDByChannelMap[channelName]
	if !exists {
		return nil, errors.New("channelName: " + channelName + " not registered")
	}
	return idmap, nil
}

func (f *subscribeFilter) Subscribe(channel, connId, userId string) {
	if connList, exists := f.connIDByChannelMap[channel]; exists {
		connList.Store(connId, true)
	} else {
		store := &sync.Map{}
		store.Store(connId, true)
		f.connIDByChannelMap[channel] = store
	}
	if connList, exists := f.connIDByUserMap[userId]; exists {
		connList.Store(connId, true)
	} else {
		store := &sync.Map{}
		store.Store(connId, true)
		f.connIDByUserMap[userId] = store
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

func (f *subscribeFilter) Unsubscribe(connId, userId string) error {
	userStore, err := f.GetMapsByUser(userId)
	if err != nil {
		return err
	}
	userStore.Delete(connId)
	if !keyExists(userStore) {
		delete(f.connIDByUserMap, userId)
	}
	for chname, store := range f.connIDByChannelMap {
		store.Delete(connId)
		if !keyExists(store) {
			delete(f.connIDByChannelMap, chname)
		}
	}

	return nil
}

func (f *subscribeFilter) GetChannelSubscriptions(channel string) map[string]bool {
	connIds := map[string]bool{}
	if connList, exists := f.connIDByChannelMap[channel]; exists {
		connList.Range(func(k, v interface{}) bool {
			connIds[k.(string)] = true
			return true
		})
	}
	return connIds
}

func (f *subscribeFilter) GetUserSubscriptions(channel string, userIds []string) map[string]bool {
	connIds := map[string]bool{}
	if connList, exists := f.connIDByChannelMap[channel]; exists {
		connList.Range(func(k, v interface{}) bool {
			connIds[k.(string)] = true
			return true
		})
	}
	userConnIds := map[string]bool{}
	for _, uid := range userIds {
		if connList, exists := f.connIDByUserMap[uid]; exists {
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
