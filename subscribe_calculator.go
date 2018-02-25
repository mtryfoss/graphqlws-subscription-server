package gss

import (
	"context"

	"github.com/graphql-go/graphql"
)

type ResolveContext struct {
	ConnectionID   string
	SubscriptionID string
	User           interface{}
	EventKey       string
	EventVal       interface{}
}

func NewResolveContext(connID, subID string, user interface{}, eKey string, eVal interface{}) *ResolveContext {
	return &ResolveContext{
		ConnectionID:   connID,
		SubscriptionID: subID,
		User:           user,
		EventKey:       eKey,
		EventVal:       eVal,
	}
}

func (r *ResolveContext) BuildContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ListenerContextKey(r.EventKey), r.EventVal)
	ctx = context.WithValue(ctx, ListenerContextKey("connID"), r.ConnectionID)
	ctx = context.WithValue(ctx, ListenerContextKey("subscriptionID"), r.SubscriptionID)
	ctx = context.WithValue(ctx, ListenerContextKey("user"), r.User)
	return ctx
}

type SubscribeCalculator interface {
	DoGraphQL(rctx *ResolveContext, query string, variables map[string]interface{}, opName string) *graphql.Result
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

func (c *subscribeCalculator) DoGraphQL(rctx *ResolveContext, query string, variables map[string]interface{}, opName string) *graphql.Result {
	return graphql.Do(graphql.Params{
		Schema:         *c.schema, // The GraphQL schema
		RequestString:  query,
		VariableValues: variables,
		OperationName:  opName,
		Context:        rctx.BuildContext(),
	})
}
