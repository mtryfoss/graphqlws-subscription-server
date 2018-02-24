package gss

import (
	"errors"

	"github.com/graphql-go/graphql"
)

type GraphQLTypeResolve interface {
	OnPayload(payload interface{}, p graphql.ResolveParams) (interface{}, error)
	OnSubscribe(l *Listener, p graphql.ResolveParams) (interface{}, error)
	OnUnsubscribe(l *Listener, p graphql.ResolveParams) (interface{}, error)
}

type GraphQLType interface {
	GetResolve() GraphQLTypeResolve
	GetType() graphql.ObjectConfig
	GetArgs() map[string]*graphql.ArgumentConfig
	FieldName() string
}

func BuildField(listener *Listener, t GraphQLType) *graphql.Field {
	args := graphql.FieldConfigArgument{}
	for name, arg := range t.GetArgs() {
		args[name] = arg
	}
	resolver := t.GetResolve()
	return &graphql.Field{
		Type: graphql.NewObject(t.GetType()),
		Args: args,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			if payload := p.Context.Value(ListenerContextKey("payload")); payload != nil { // payload exists
				return resolver.OnPayload(payload, p)
			}
			if s := p.Context.Value(ListenerContextKey("onSubscribe")); s != nil { // AddSubscription called
				return resolver.OnSubscribe(listener, p)
			}
			if s := p.Context.Value(ListenerContextKey("onUnsubscribe")); s != nil { // removeSubscription called
				return resolver.OnUnsubscribe(listener, p)
			}
			return nil, errors.New("no payload exists")
		},
	}
}
