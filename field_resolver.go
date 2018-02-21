package graphqlws_subscription_server

import (
	"context"
	"errors"

	"github.com/graphql-go/graphql"
)

type GraphQLType interface {
	OnPayload(interface{}, context.Context) (interface{}, error)
	OnSubscribe(context.Context, *Listener) (interface{}, error)
	OnUnsubscribe(context.Context, *Listener) (interface{}, error)
	GetField(*Listener) *graphql.Field
	FieldName() string
}

type GraphQLTypeImpl struct {
	GraphQLType
	fieldName string
}

func (t *GraphQLTypeImpl) FieldName() string {
	return t.fieldName
}

func (t *GraphQLTypeImpl) GetResolve(l *Listener) func(p graphql.ResolveParams) (interface{}, error) {
	return func(p graphql.ResolveParams) (interface{}, error) {
		if payload := p.Context.Value("payload"); payload != nil { // payload exists
			return t.OnPayload(payload, p.Context)
		}
		if s := p.Context.Value("onSubscribe"); s != nil { // AddSubscription called
			return t.OnSubscribe(p.Context, l)
		}
		if s := p.Context.Value("onUnsubscribe"); s != nil { // removeSubscription called
			return t.OnUnsubscribe(p.Context, l)
		}
		return nil, errors.New("no payload exists")
	}
}
