package gss

import (
	"context"

	"github.com/functionalfoundry/graphqlws"
	"github.com/graphql-go/graphql"
)

type SubscribeCalculator interface {
	Do(ctx context.Context, query string, variables map[string]interface{}, opName string) *graphql.Result
}

type subscribeCalculator struct {
	SubscribeCalculator
	schema *graphql.Schema
}

func NewSubscribeCalculator(schema *graphql.Schema) *subscribeCalculator {
	return &subscribeCalculator{
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

func (c *subscribeCalculator) Do(ctx context.Context, query string, variables map[string]interface{}, opName string) *graphql.Result {
	return graphql.Do(graphql.Params{
		Schema:         *c.schema, // The GraphQL schema
		RequestString:  query,
		VariableValues: variables,
		OperationName:  opName,
		Context:        ctx,
	})
}
