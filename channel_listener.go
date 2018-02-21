package graphqlws_subscription_server

import (
	"context"
	"sync"

	"github.com/functionalfoundry/graphqlws"
	"github.com/graphql-go/graphql"
)

type ConnectionsByID map[string]bool

type Listener struct {
	graphqlws.SubscriptionManager
	manager            *graphqlws.SubscriptionManager
	schema             *graphql.Schema
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

func (l *Listener) BuildManager(schema *graphql.Schema) {
	l.schema = schema
	m := graphqlws.NewSubscriptionManager(schema)
	l.manager = &m
}

func bldCtx(flg string, conn graphqlws.Connection) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, flg, true)
	ctx = context.WithValue(ctx, "connID", conn.ID())
	ctx = context.WithValue(ctx, "user", conn.User())
	return ctx
}

func (l *Listener) doGraphQL(ctx context.Context, s *graphqlws.Subscription) *graphql.Result {
	return graphql.Do(graphql.Params{
		Schema:         *l.schema, // The GraphQL schema
		RequestString:  s.Query,
		VariableValues: s.Variables,
		OperationName:  s.OperationName,
		Context:        ctx,
	})
}

func (l *Listener) AddSubscription(conn graphqlws.Connection, s *graphqlws.Subscription) []error {
	result := l.doGraphQL(bldCtx("onSubscribe", conn), s)

	if result.HasErrors() {
		return graphqlws.ErrorsFromGraphQLErrors(result.Errors)
	}

	errs := (*l.manager).AddSubscription(conn, s)
	if errs != nil {
		return errs
	}

	return nil
}

func (l *Listener) RemoveSubscription(conn graphqlws.Connection, s *graphqlws.Subscription) {
	l.doGraphQL(bldCtx("onUnsubscribe", conn), s)
	(*l.manager).RemoveSubscription(conn, s)
}

func (l *Listener) RemoveSubscriptions(conn graphqlws.Connection) {
	ctx := bldCtx("onUnsubscribe", conn)
	for _, subscription := range l.Subscriptions()[conn] {
		l.doGraphQL(ctx, subscription)
	}
	(*l.manager).RemoveSubscriptions(conn)
}

func (l *Listener) Subscriptions() graphqlws.Subscriptions {
	return (*l.manager).Subscriptions()
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
