package gss

import (
	"context"

	"github.com/functionalfoundry/graphqlws"
	"github.com/graphql-go/graphql"
)

type ChannelExecutor interface {
	DoGraphQL(context.Context, *graphqlws.Subscription) *graphql.Result
}

type channelExecuter struct {
	ChannelExecutor
	schema *graphql.Schema
}

type ListenerContextKey string

func NewChannelExecutor(schema *graphql.Schema) ChannelExecutor {
	return &channelExecuter{
		schema: schema,
	}
}

func buildCtx(eventName string, eventVal interface{}, conn graphqlws.Connection) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ListenerContextKey(eventName), eventVal)
	ctx = context.WithValue(ctx, ListenerContextKey("connID"), conn.ID())
	ctx = context.WithValue(ctx, ListenerContextKey("user"), conn.User())
	return ctx
}

func (e *channelExecuter) DoGraphQL(ctx context.Context, s *graphqlws.Subscription) *graphql.Result {
	return graphql.Do(graphql.Params{
		Schema:         *e.schema, // The GraphQL schema
		RequestString:  s.Query,
		VariableValues: s.Variables,
		OperationName:  s.OperationName,
		Context:        ctx,
	})
}
