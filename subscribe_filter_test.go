package gss

import (
	"sync"
	"testing"
)

func TestChannelManagerSubscribeAndUnsubscribe(t *testing.T) {
	f := NewSubscribeFilter()
	chanName1 := "foo"
	subscriptionID1 := "sub1"
	connID1 := "conn1"
	userID1 := "hoge"
	f.Subscribe(chanName1, subscriptionID1, connID1, userID1)
	idmap, err := f.GetMapsByUser("fuga")
	if idmap != nil {
		t.Error("wrong user map exists")
	}
	if err == nil {
		t.Error("error should exists")
	}
	if err.Error() != "userID: fuga not registered" {
		t.Error("unexpected error message")
	}
	idmap, err = f.GetMapsByUser(userID1)
	if idmap == nil {
		t.Error("idmap should exists")
	}
	if err != nil {
		t.Error("error should not exists")
	}
	if _, ok := idmap.Load(subscriptionID1); !ok {
		t.Error("subscriptionID1 not registered")
	}

	var connections ConnIDBySubscriptionID
	var connectionsB map[string]bool

	subscriptionID2 := "sub2"
	connID2 := "conn2"
	userID2 := "fuga"
	f.Subscribe(chanName1, subscriptionID2, connID2, userID2)

	chanName2 := "bar"
	subscriptionID3 := "conn3"
	f.Subscribe(chanName2, subscriptionID3, connID1, userID1)

	idmap, _ = f.GetMapsByUser(userID1)
	connectionsB = getConnsFromSyncMapB(idmap)
	if len(connectionsB) != 2 {
		t.Error("userID1 connections not enough")
	}
	if _, exists := connectionsB[subscriptionID1]; !exists {
		t.Error("subscriptionID1 not registered")
	}
	if _, exists := connectionsB[subscriptionID3]; !exists {
		t.Error("subscriptionID3 not registered")
	}

	idmap, err = f.GetMapsByChannel("baz")
	if idmap != nil {
		t.Error("channel: baz is not registered")
	}
	if err == nil {
		t.Error("error should exists")
	}
	if err.Error() != "channelName: baz not registered" {
		t.Error("unexpected error message")
	}
	idmap, err = f.GetMapsByChannel(chanName1)
	if idmap == nil {
		t.Error("idmap should exists")
	}
	if err != nil {
		t.Error("error should not exists")
	}
	connections = getConnsFromSyncMap(idmap)
	if len(connections) != 2 {
		t.Error("chanName1 connections not enough")
	}
	if _, exists := connections[subscriptionID1]; !exists {
		t.Error("subscriptionID1 should exists")
	}
	if _, exists := connections[subscriptionID2]; !exists {
		t.Error("subscriptionID2 should exists")
	}
	idmap, _ = f.GetMapsByChannel(chanName2)
	connections = getConnsFromSyncMap(idmap)
	if len(connections) != 1 {
		t.Error("chanName2 connections not enough")
	}
	if _, exists := connections[subscriptionID3]; !exists {
		t.Error("subscriptionID3 should exists")
	}

	err = f.Unsubscribe(subscriptionID1, userID1)
	if err != nil {
		t.Error("error should not exists")
	}

	idmap, _ = f.GetMapsByUser(userID1)
	connectionsB = getConnsFromSyncMapB(idmap)
	if len(connectionsB) != 1 {
		t.Error("rest connections is 1")
	}
	if _, exists := connectionsB[subscriptionID3]; !exists {
		t.Error("subscriptionID3 should exists")
	}
	idmap, _ = f.GetMapsByUser(userID2)
	connectionsB = getConnsFromSyncMapB(idmap)
	if len(connections) != 1 {
		t.Error("rest connections is 1")
	}
	if _, exists := connectionsB[subscriptionID2]; !exists {
		t.Error("subscriptionID2 should exists")
	}

	idmap, _ = f.GetMapsByChannel(chanName1)
	connections = getConnsFromSyncMap(idmap)
	if len(connections) != 1 {
		t.Error("rest connections is 1")
	}
	if _, exists := connections[subscriptionID2]; !exists {
		t.Error("subscriptionID2 should exists")
	}

	idmap, _ = f.GetMapsByChannel(chanName2)
	connections = getConnsFromSyncMap(idmap)
	if len(connections) != 1 {
		t.Error("rest connections is 1")
	}
	if _, exists := connections[subscriptionID3]; !exists {
		t.Error("subscriptionID2 should exists")
	}

	f.Unsubscribe(subscriptionID2, userID2)

	idmap, err = f.GetMapsByChannel(chanName1)
	if idmap != nil {
		t.Error("chanName1 map should removed")
	}
	if err == nil {
		t.Error("error should exists")
	}
	if err.Error() != "channelName: foo not registered" {
		t.Error("unexpected error message")
	}

	idmap, err = f.GetMapsByUser(userID2)
	if idmap != nil {
		t.Error("userID2 map should removed")
	}
	if err == nil {
		t.Error("error should exists")
	}
	if err.Error() != "userID: fuga not registered" {
		t.Error("unexpected error message")
	}

	err = f.Unsubscribe("111111", userID2)
	if err == nil {
		t.Error("error should exists")
	}
	if err.Error() != "userID: fuga not registered" {
		t.Error("unexpected error message")
	}
}

func TestGetSubscriptionsByChannelManager(t *testing.T) {
	f := NewSubscribeFilter()
	chName1 := "foo"
	chName2 := "bar"
	subs := f.GetChannelSubscriptions(chName1)
	if len(subs) != 0 {
		t.Error("subscriptions count should be 0")
	}
	f.Subscribe(chName1, "conn1", "sub1", "user1")
	f.Subscribe(chName1, "conn2", "sub1", "user2")
	f.Subscribe(chName2, "conn3", "sub2", "user3")
	f.Subscribe(chName1, "conn4", "sub1", "user4")
	f.Subscribe(chName1, "conn5", "sub2", "user4")

	subs = f.GetChannelSubscriptions(chName1)
	if len(subs) != 4 {
		t.Error("subscriptions count should be 4")
	}
	for _, c := range []string{"conn1", "conn2", "conn4", "conn5"} {
		if _, exists := subs[c]; !exists {
			t.Error("connection: " + c + " not found")
		}
	}

	subs = f.GetUserSubscriptions(chName1, []string{"user4"})
	if len(subs) != 2 {
		t.Error("subscriptions count should be 2")
	}
	if _, exists := subs["conn4"]; !exists {
		t.Error("conn4 should exists")
	}
	if _, exists := subs["conn5"]; !exists {
		t.Error("conn5 should exists")
	}

	subs = f.GetUserSubscriptions(chName1, []string{"user3"})
	if len(subs) != 0 {
		t.Error("subscriptions count should be 0")
	}

}

func getConnsFromSyncMap(m *sync.Map) ConnIDBySubscriptionID {
	connections := ConnIDBySubscriptionID{}
	m.Range(func(k, v interface{}) bool {
		connections[k.(string)] = v.(string)
		return true
	})
	return connections
}

func getConnsFromSyncMapB(m *sync.Map) map[string]bool {
	connections := map[string]bool{}
	m.Range(func(k, v interface{}) bool {
		connections[k.(string)] = true
		return true
	})
	return connections
}
