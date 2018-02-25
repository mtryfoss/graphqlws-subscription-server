package gss

import (
	"errors"
	"sync"
)

type ConnIDBySubscriptionID map[string]string

type SubscribeFilter interface {
	Subscribe(channelName, subscriptionID, connID, userID string)
	Unsubscribe(subscriptionID, userID string) error
	GetChannelSubscriptionIDs(channelName string) ConnIDBySubscriptionID
	GetUserSubscriptionIDs(channelName string, userIDs []string) ConnIDBySubscriptionID
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

func (f *subscribeFilter) Subscribe(channel, subscriptionID, connID, userId string) {
	if connList, exists := f.subscriptionIDByChannelMap[channel]; exists {
		connList.Store(subscriptionID, connID)
	} else {
		store := &sync.Map{}
		store.Store(subscriptionID, connID)
		f.subscriptionIDByChannelMap[channel] = store
	}
	if connList, exists := f.subscriptionIDByUserMap[userId]; exists {
		connList.Store(subscriptionID, connID)
	} else {
		store := &sync.Map{}
		store.Store(subscriptionID, connID)
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

func (f *subscribeFilter) GetChannelSubscriptions(channel string) ConnIDBySubscriptionID {
	subscriptionIDs := ConnIDBySubscriptionID{}
	if connList, exists := f.subscriptionIDByChannelMap[channel]; exists {
		connList.Range(func(k, v interface{}) bool {
			subscriptionIDs[k.(string)] = v.(string)
			return true
		})
	}
	return subscriptionIDs
}

func (f *subscribeFilter) GetUserSubscriptions(channel string, userIds []string) ConnIDBySubscriptionID {
	subscriptionIDs := ConnIDBySubscriptionID{}
	if connList, exists := f.subscriptionIDByChannelMap[channel]; exists {
		connList.Range(func(k, v interface{}) bool {
			subscriptionIDs[k.(string)] = v.(string)
			return true
		})
	}
	usersubscriptionIDs := ConnIDBySubscriptionID{}
	for _, uid := range userIds {
		if connList, exists := f.subscriptionIDByUserMap[uid]; exists {
			connList.Range(func(k, v interface{}) bool {
				subscriptionID := k.(string)
				if v, exists := subscriptionIDs[subscriptionID]; exists {
					usersubscriptionIDs[subscriptionID] = v
				}
				return true
			})
		}
	}
	return usersubscriptionIDs
}
