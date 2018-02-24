package gss

import (
	"sync"
	"testing"
)

func TestChannelManagerSubscribeAndUnsubscribe(t *testing.T) {
	m := NewSubscribeFilter()
	chanName1 := "foo"
	connID1 := "conn1"
	userID1 := "hoge"
	m.Subscribe(chanName1, connID1, userID1)
	idmap, err := m.GetMapsByUser("fuga")
	if idmap != nil {
		t.Error("wrong user map exists")
	}
	if err == nil {
		t.Error("error should exists")
	}
	if err.Error() != "userID: fuga not registered" {
		t.Error("unexpected error message")
	}
	idmap, err = m.GetMapsByUser(userID1)
	if idmap == nil {
		t.Error("idmap should exists")
	}
	if err != nil {
		t.Error("error should not exists")
	}
	if _, ok := idmap.Load(connID1); !ok {
		t.Error("connID1 not registered")
	}

	connID2 := "conn2"
	userID2 := "fuga"
	m.Subscribe(chanName1, connID2, userID2)

	chanName2 := "bar"
	connID3 := "conn3"
	m.Subscribe(chanName2, connID3, userID1)

	idmap, _ = m.GetMapsByUser(userID1)
	connections := getConnsFromSyncMap(idmap)
	if len(connections) != 2 {
		t.Error("userID1 connections not enough")
	}
	if _, exists := connections[connID1]; !exists {
		t.Error("connID1 not registered")
	}
	if _, exists := connections[connID3]; !exists {
		t.Error("connID3 not registered")
	}

	idmap, err = m.GetMapsByChannel("baz")
	if idmap != nil {
		t.Error("channel: baz is not registered")
	}
	if err == nil {
		t.Error("error should exists")
	}
	if err.Error() != "channelName: baz not registered" {
		t.Error("unexpected error message")
	}
	idmap, err = m.GetMapsByChannel(chanName1)
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
	if _, exists := connections[connID1]; !exists {
		t.Error("connID1 should exists")
	}
	if _, exists := connections[connID2]; !exists {
		t.Error("connID2 should exists")
	}
	idmap, _ = m.GetMapsByChannel(chanName2)
	connections = getConnsFromSyncMap(idmap)
	if len(connections) != 1 {
		t.Error("chanName2 connections not enough")
	}
	if _, exists := connections[connID3]; !exists {
		t.Error("connID3 should exists")
	}

	err = m.Unsubscribe(connID1, userID1)
	if err != nil {
		t.Error("error should not exists")
	}

	idmap, _ = m.GetMapsByUser(userID1)
	connections = getConnsFromSyncMap(idmap)
	if len(connections) != 1 {
		t.Error("rest connections is 1")
	}
	if _, exists := connections[connID3]; !exists {
		t.Error("connID3 should exists")
	}
	idmap, _ = m.GetMapsByUser(userID2)
	connections = getConnsFromSyncMap(idmap)
	if len(connections) != 1 {
		t.Error("rest connections is 1")
	}
	if _, exists := connections[connID2]; !exists {
		t.Error("connID2 should exists")
	}

	idmap, _ = m.GetMapsByChannel(chanName1)
	connections = getConnsFromSyncMap(idmap)
	if len(connections) != 1 {
		t.Error("rest connections is 1")
	}
	if _, exists := connections[connID2]; !exists {
		t.Error("connID2 should exists")
	}

	idmap, _ = m.GetMapsByChannel(chanName2)
	connections = getConnsFromSyncMap(idmap)
	if len(connections) != 1 {
		t.Error("rest connections is 1")
	}
	if _, exists := connections[connID3]; !exists {
		t.Error("connID2 should exists")
	}

	m.Unsubscribe(connID2, userID2)

	idmap, err = m.GetMapsByChannel(chanName1)
	if idmap != nil {
		t.Error("chanName1 map should removed")
	}
	if err == nil {
		t.Error("error should exists")
	}
	if err.Error() != "channelName: foo not registered" {
		t.Error("unexpected error message")
	}

	idmap, err = m.GetMapsByUser(userID2)
	if idmap != nil {
		t.Error("userID2 map should removed")
	}
	if err == nil {
		t.Error("error should exists")
	}
	if err.Error() != "userID: fuga not registered" {
		t.Error("unexpected error message")
	}

	err = m.Unsubscribe("111111", userID2)
	if err == nil {
		t.Error("error should exists")
	}
	if err.Error() != "userID: fuga not registered" {
		t.Error("unexpected error message")
	}
}

func TestGetSubscriptionsByChannelManager(t *testing.T) {
	m := NewSubscribeFilter()
	chName1 := "foo"
	chName2 := "bar"
	subs := m.GetChannelSubscriptions(chName1)
	if len(subs) != 0 {
		t.Error("subscriptions count should be 0")
	}
	m.Subscribe(chName1, "conn1", "user1")
	m.Subscribe(chName1, "conn2", "user2")
	m.Subscribe(chName2, "conn3", "user3")
	m.Subscribe(chName1, "conn4", "user4")
	m.Subscribe(chName1, "conn5", "user4")

	subs = m.GetChannelSubscriptions(chName1)
	if len(subs) != 4 {
		t.Error("subscriptions count should be 4")
	}
	for _, c := range []string{"conn1", "conn2", "conn4", "conn5"} {
		if _, exists := subs[c]; !exists {
			t.Error("connection: " + c + " not found")
		}
	}

	subs = m.GetUserSubscriptions(chName1, []string{"user4"})
	if len(subs) != 2 {
		t.Error("subscriptions count should be 2")
	}
	if _, exists := subs["conn4"]; !exists {
		t.Error("conn4 should exists")
	}
	if _, exists := subs["conn5"]; !exists {
		t.Error("conn5 should exists")
	}

	subs = m.GetUserSubscriptions(chName1, []string{"user3"})
	if len(subs) != 0 {
		t.Error("subscriptions count should be 0")
	}

}

func getConnsFromSyncMap(m *sync.Map) map[string]bool {
	connections := map[string]bool{}
	m.Range(func(k, v interface{}) bool {
		connections[k.(string)] = true
		return true
	})
	return connections
}
