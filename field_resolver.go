package graphqlws_subscription_server

import (
	"errors"

	"github.com/graphql-go/graphql"
)

type GraphQLType interface {
	OnPayload(interface{}, graphql.ResolveParams) (interface{}, error)
	OnSubscribe(graphql.ResolveParams, *Listener) (interface{}, error)
	OnUnsubscribe(graphql.ResolveParams, *Listener) (interface{}, error)
	GetField(*Listener) *graphql.Field
	GetType() graphql.ObjectConfig
	GetArgs() map[string]*graphql.ArgumentConfig
	FieldName() string
}

type GraphQLTypeImpl struct {
	GraphQLType
	fieldName string
}

func (t *GraphQLTypeImpl) FieldName() string {
	return t.fieldName
}

func (t *GraphQLTypeImpl) GetField(listener *Listener) *graphql.Field {
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
				return t.OnSubscribe(p, listener)
			}
			if s := p.Context.Value(ListenerContextKey("onUnsubscribe")); s != nil { // removeSubscription called
				return t.OnUnsubscribe(p, listener)
			}
			return nil, errors.New("no payload exists")
		},
	}
}
