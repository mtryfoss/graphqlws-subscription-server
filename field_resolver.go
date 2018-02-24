package gss

import (
	"errors"

	"github.com/graphql-go/graphql"
)

type GraphQLType interface {
	OnPayload(payload interface{}, p graphql.ResolveParams) (interface{}, error)
	OnSubscribe(l *Listener, p graphql.ResolveParams) (interface{}, error)
	OnUnsubscribe(l *Listener, p graphql.ResolveParams) (interface{}, error)
	GetType() graphql.ObjectConfig
	GetArgs() map[string]*graphql.ArgumentConfig
	FieldName() string
}

func buildField(listener *Listener, t GraphQLType) *graphql.Field {
	args := graphql.FieldConfigArgument{}
	for name, arg := range t.GetArgs() {
		args[name] = arg
	}
	return &graphql.Field{
		Type: graphql.NewObject(t.GetType()),
		Args: args,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			if payload := p.Context.Value(ListenerContextKey("payload")); payload != nil { // payload exists
				return t.OnPayload(payload, p)
			}
			if s := p.Context.Value(ListenerContextKey("onSubscribe")); s != nil { // AddSubscription called
				return t.OnSubscribe(listener, p)
			}
			if s := p.Context.Value(ListenerContextKey("onUnsubscribe")); s != nil { // removeSubscription called
				return t.OnUnsubscribe(listener, p)
			}
			return nil, errors.New("no payload exists")
		},
	}
}
