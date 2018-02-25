package gss

import (
	"testing"

	"github.com/graphql-go/graphql"
)

func TestNewSubscribeService(t *testing.T) {
	user := map[string]string{}
	user["foo"] = "world"

	// Schema (via https://github.com/graphql-go/graphql )
	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				v, _ := p.Context.Value(GraphQLContextKey("user")).(map[string]string)["foo"]
				return v, nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, _ := graphql.NewSchema(schemaConfig)

	subChan := make(chan *SubscribeEvent, 1)
	unsubChan := make(chan *UnsubscribeEvent, 1)

	s := NewSubscribeService(&schema, 10, subChan, unsubChan)
	s.SubscribeFilter()
	s.SubscribeCalculator()
	s.GetNotifierChan()
	s.Subscriptions()
}
