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
	connIDByUserMap    map[string]*sync.Map
	connIDByChannelMap map[string]*sync.Map
	dummyLabel         string
}

func NewListener(dummy string) *Listener {
	return &Listener{
		connIDByUserMap:    map[string]*sync.Map{},
		connIDByChannelMap: map[string]*sync.Map{},
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
	if connList, exists := l.connIDByChannelMap[channel]; exists {
		connList.Store(connId, true)
	} else {
		store := &sync.Map{}
		store.Store(connId, true)
		l.connIDByChannelMap[channel] = store
	}
	if userId == l.dummyLabel {
		return
	}
	if connList, exists := l.connIDByUserMap[userId]; exists {
		connList.Store(connId, true)
	} else {
		store := &sync.Map{}
		store.Store(connId, true)
		l.connIDByUserMap[userId] = store
	}
}

func (l *Listener) Unsubscribe(connId, userId string) {
	connIds := []string{connId}
	if store, exists := l.connIDByUserMap[userId]; exists {
		store.Range(func(k, v interface{}) bool {
			connIds = append(connIds, k.(string))
			return true
		})
		delete(l.connIDByUserMap, userId)
	}
	for _, connList := range l.connIDByChannelMap {
		for _, cid := range connIds {
			connList.Delete(cid)
		}
	}
}

func (l *Listener) GetChannelSubscriptions(channel string) graphqlws.Subscriptions {
	subscriptions := graphqlws.Subscriptions{}
	if connList, exists := l.connIDByChannelMap[channel]; exists {
		for conn, s := range l.Subscriptions() {
			if _, ok := connList.Load(conn.ID()); ok {
				subscriptions[conn] = s
			}
		}
	}
	return subscriptions
}

func (l *Listener) GetUserSubscriptions(userIds []string) graphqlws.Subscriptions {
	subscriptions := graphqlws.Subscriptions{}
	connIds := map[string]bool{}
	for _, uid := range userIds {
		if connList, exists := l.connIDByUserMap[uid]; exists {
			connList.Range(func(k, v interface{}) bool {
				connIds[k.(string)] = true
				return true
			})
		}
	}
	for conn, s := range l.Subscriptions() {
		if _, exists := connIds[conn.ID()]; exists {
			subscriptions[conn] = s
		}
	}
	return subscriptions
}
