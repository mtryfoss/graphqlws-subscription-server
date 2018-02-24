package gss

import (
	"context"
	"sync"

	"github.com/functionalfoundry/graphqlws"
	"github.com/graphql-go/graphql"
)

type Listener struct {
	graphqlws.SubscriptionManager
	ms         *graphqlws.SubscriptionManager
	mc         ChannelManager
	me         ChannelExecutor
	notifyChan chan *RequestData
	subChan    chan *SubscribeEvent
	unsubChan  chan *UnsubscribeEvent
}

func NewListener(schema *graphql.Schema, handleCount uint, subChan chan *SubscribeEvent, unsubChan chan *UnsubscribeEvent) *Listener {
	ms := graphqlws.NewSubscriptionManager(schema)
	mc := NewChannelManager()
	me := NewChannelExecutor(schema)
	return &Listener{
		ms:         &ms,
		mc:         mc,
		me:         me,
		notifyChan: make(chan *RequestData, handleCount),
		subChan:    subChan,
		unsubChan:  unsubChan,
	}
}

func (l *Listener) AddSubscription(conn graphqlws.Connection, s *graphqlws.Subscription) []error {
	l.me.DoGraphQL(buildCtx("onSubscribe", true, conn), s)
	errs := (*l.ms).AddSubscription(conn, s)
	if errs != nil {
		return errs
	}

	return nil
}

func (l *Listener) RemoveSubscription(conn graphqlws.Connection, s *graphqlws.Subscription) {
	l.me.DoGraphQL(buildCtx("onUnsubscribe", true, conn), s)
	(*l.ms).RemoveSubscription(conn, s)
}

func (l *Listener) RemoveSubscriptions(conn graphqlws.Connection) {
	ctx := buildCtx("onUnsubscribe", true, conn)
	for _, subscription := range l.Subscriptions()[conn] {
		l.me.DoGraphQL(ctx, subscription)
	}
	(*l.ms).RemoveSubscriptions(conn)
}

func (l *Listener) Subscriptions() graphqlws.Subscriptions {
	return (*l.ms).Subscriptions()
}

func (l *Listener) Publish(connIds map[string]bool, payload interface{}) {
	for conn, _ := range l.Subscriptions() {
		if _, exists := connIds[conn.ID()]; exists {
			for _, subscription := range l.Subscriptions()[conn] {
				res := l.me.DoGraphQL(buildCtx("payload", payload, conn), subscription)
				d := &graphqlws.DataMessagePayload{
					Data: res.Data,
				}
				if res.HasErrors() {
					d.Errors = graphqlws.ErrorsFromGraphQLErrors(res.Errors)
				}
				subscription.SendData(d)
			}
		}
	}
}

func (l *Listener) ChannelManager() ChannelManager {
	return l.mc
}

func (l *Listener) ChannelExecuter() ChannelExecutor {
	return l.me
}

func (l *Listener) GetNotifierChan() chan *RequestData {
	return l.notifyChan
}

func (l *Listener) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	chanManager := l.ChannelManager()
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case data := <-l.notifyChan:
				var connIds map[string]bool
				if len(data.Users) > 0 {
					connIds = chanManager.GetUserSubscriptions(data.Channel, data.Users)
				} else {
					connIds = chanManager.GetChannelSubscriptions(data.Channel)
				}
				if len(connIds) > 0 {
					go l.Publish(connIds, data.Payload)
				}
			case data := <-l.subChan:
				l.mc.Subscribe(data.Channel, data.ConnID, data.User)
			case data := <-l.unsubChan:
				l.mc.Unsubscribe(data.ConnID, data.User)
			}
		}
	}()
}
