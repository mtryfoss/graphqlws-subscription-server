package gss

import (
	"errors"
	"sync"
)

type SubscribeFilter interface {
	Subscribe(channelName, subscriptionID, userID string)
	Unsubscribe(subscriptionID, userID string) error
	GetChannelSubscriptions(channelName string) map[string]bool
	GetUserSubscriptions(channelName string, userIDs []string) map[string]bool
}

type subscribeFilter struct {
	SubscribeFilter
	subscriptionIDByUserMap    map[string]*sync.Map
	subscriptionIDByChannelMap map[string]*sync.Map
}

func NewSubscribeFilter() *subscribeFilter {
	return &subscribeFilter{
		subscriptionIDByUserMap:    map[string]*sync.Map{},
		subscriptionIDByChannelMap: map[string]*sync.Map{},
	}
}

func (f *subscribeFilter) GetMapsByUser(userID string) (*sync.Map, error) {
	idmap, exists := f.subscriptionIDByUserMap[userID]
	if !exists {
		return nil, errors.New("userID: " + userID + " not registered")
	}
	return idmap, nil
}

func (f *subscribeFilter) GetMapsByChannel(channelName string) (*sync.Map, error) {
	idmap, exists := f.subscriptionIDByChannelMap[channelName]
	if !exists {
		return nil, errors.New("channelName: " + channelName + " not registered")
	}
	return idmap, nil
}

func (f *subscribeFilter) Subscribe(channel, subscriptionID, userId string) {
	if connList, exists := f.subscriptionIDByChannelMap[channel]; exists {
		connList.Store(subscriptionID, true)
	} else {
		store := &sync.Map{}
		store.Store(subscriptionID, true)
		f.subscriptionIDByChannelMap[channel] = store
	}
	if connList, exists := f.subscriptionIDByUserMap[userId]; exists {
		connList.Store(subscriptionID, true)
	} else {
		store := &sync.Map{}
		store.Store(subscriptionID, true)
		f.subscriptionIDByUserMap[userId] = store
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

func (f *subscribeFilter) Unsubscribe(subscriptionID, userId string) error {
	userStore, err := f.GetMapsByUser(userId)
	if err != nil {
		return err
	}
	userStore.Delete(subscriptionID)
	if !keyExists(userStore) {
		delete(f.subscriptionIDByUserMap, userId)
	}
	for chname, store := range f.subscriptionIDByChannelMap {
		store.Delete(subscriptionID)
		if !keyExists(store) {
			delete(f.subscriptionIDByChannelMap, chname)
		}
	}

	return nil
}

func (f *subscribeFilter) GetChannelSubscriptions(channel string) map[string]bool {
	subscriptionIDs := map[string]bool{}
	if connList, exists := f.subscriptionIDByChannelMap[channel]; exists {
		connList.Range(func(k, v interface{}) bool {
			subscriptionIDs[k.(string)] = true
			return true
		})
	}
	return subscriptionIDs
}

func (f *subscribeFilter) GetUserSubscriptions(channel string, userIds []string) map[string]bool {
	subscriptionIDs := map[string]bool{}
	if connList, exists := f.subscriptionIDByChannelMap[channel]; exists {
		connList.Range(func(k, v interface{}) bool {
			subscriptionIDs[k.(string)] = true
			return true
		})
	}
	usersubscriptionIDs := map[string]bool{}
	for _, uid := range userIds {
		if connList, exists := f.subscriptionIDByUserMap[uid]; exists {
			connList.Range(func(k, v interface{}) bool {
				subscriptionID := k.(string)
				if _, exists := subscriptionIDs[subscriptionID]; exists {
					usersubscriptionIDs[subscriptionID] = true
				}
				return true
			})
		}
	}
	return usersubscriptionIDs
}
