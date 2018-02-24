package gss

import (
	"context"

	"github.com/functionalfoundry/graphqlws"
	"github.com/graphql-go/graphql"
)

type Listener struct {
	graphqlws.SubscriptionManager
	ms     *graphqlws.SubscriptionManager
	mc     *ChannelManager
	schema *graphql.Schema
}

type ListenerContextKey string

func NewListener(schema *graphql.Schema) *Listener {
	ms := graphqlws.NewSubscriptionManager(schema)
	mc := NewChannelManager()
	return &Listener{
		schema: schema,
		ms:     &ms,
		mc:     mc,
	}
}

func buildCtx(eventName string, eventVal interface{}, conn graphqlws.Connection) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ListenerContextKey(eventName), eventVal)
	ctx = context.WithValue(ctx, ListenerContextKey("connID"), conn.ID())
	ctx = context.WithValue(ctx, ListenerContextKey("user"), conn.User())
	return ctx
}

func (l *Listener) DoGraphQL(ctx context.Context, s *graphqlws.Subscription) *graphql.Result {
	return graphql.Do(graphql.Params{
		Schema:         *l.schema, // The GraphQL schema
		RequestString:  s.Query,
		VariableValues: s.Variables,
		OperationName:  s.OperationName,
		Context:        ctx,
	})
}

func (l *Listener) AddSubscription(conn graphqlws.Connection, s *graphqlws.Subscription) []error {
	result := l.DoGraphQL(buildCtx("onSubscribe", true, conn), s)

	if result.HasErrors() {
		return graphqlws.ErrorsFromGraphQLErrors(result.Errors)
	}

	errs := (*l.ms).AddSubscription(conn, s)
	if errs != nil {
		return errs
	}

	return nil
}

func (l *Listener) RemoveSubscription(conn graphqlws.Connection, s *graphqlws.Subscription) {
	l.DoGraphQL(buildCtx("onUnsubscribe", true, conn), s)
	(*l.ms).RemoveSubscription(conn, s)
}

func (l *Listener) RemoveSubscriptions(conn graphqlws.Connection) {
	ctx := buildCtx("onUnsubscribe", true, conn)
	for _, subscription := range l.Subscriptions()[conn] {
		l.DoGraphQL(ctx, subscription)
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
				res := l.DoGraphQL(buildCtx("payload", payload, conn), subscription)
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

func (l *Listener) ChannelManager() *ChannelManager {
	return l.mc
}
