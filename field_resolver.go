package gss

import (
	"errors"

	"github.com/graphql-go/graphql"
)

type GraphQLResolve interface {
	OnPayload(payload interface{}, p graphql.ResolveParams) (interface{}, error)
	OnSubscribe(p graphql.ResolveParams) (interface{}, error)
	OnUnsubscribe(p graphql.ResolveParams) (interface{}, error)
}

func BuildResolve(r GraphQLResolve) func(p graphql.ResolveParams) (interface{}, error) {
	return func(p graphql.ResolveParams) (interface{}, error) {
		if payload := p.Context.Value(GraphQLContextKey("payload")); payload != nil { // payload exists
			return r.OnPayload(payload, p)
		}
		if s := p.Context.Value(GraphQLContextKey("onSubscribe")); s != nil { // AddSubscription called
			return r.OnSubscribe(p)
		}
		if s := p.Context.Value(GraphQLContextKey("onUnsubscribe")); s != nil { // removeSubscription called
			return r.OnUnsubscribe(p)
		}
		return nil, errors.New("no payload exists")
	}
}
