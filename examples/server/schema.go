package main

import (
	"github.com/graphql-go/graphql"
	gss "github.com/taiyoh/graphqlws-subscription-server"
)

// this is sample code.
func LoadSchema(listener *gss.Listener) (graphql.Schema, error) {
	fields := graphql.Fields{}
	resolvers: = []TypeResolver{NewCommentResolver()}
	for _, resolver := range resolvers {
		name, field := t.GetField(listener)
		fields[name] = field
	}
	return graphql.NewSchema(graphql.SchemaConfig{
		Subscription: graphql.NewObject(
			graphql.ObjectConfig{
				Name: "RootSubscription",
				Fields: fields,
			}
		),
	})
}
