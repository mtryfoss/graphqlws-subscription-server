package gss

import (
	"testing"

	"github.com/graphql-go/graphql"
)

func TestNewResolveContext(t *testing.T) {
	user := map[string]string{}
	user["foo"] = "bar"
	rctx := NewResolveContext("connID1", "subscriptionID1", user, "ek1", "ev1")
	ctx := rctx.BuildContext()

	if ctx.Value(GraphQLContextKey("connID")).(string) != "connID1" {
		t.Error("ConnectionID is expected as connID1")
	}
	if ctx.Value(GraphQLContextKey("subscriptionID")).(string) != "subscriptionID1" {
		t.Error("SubscriptionID is expected as subscriptionID1")
	}
	if ctx.Value(GraphQLContextKey("user")) == nil {
		t.Error("User should exists")
	}
	if ctx.Value(GraphQLContextKey("ek1")).(string) != "ev1" {
		t.Error("EventKey is expected as ek1 and EventVal is expected as ev1")
	}
}

func TestDoGraphQL(t *testing.T) {
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

	c := NewSubscribeCalculator(&schema)

	rctx := NewResolveContext("connID1", "subscriptionID1", user, "ek1", "ev1")

	// Query
	query := `
		{
			hello
		}
	`
	v := map[string]interface{}{}
	res := c.DoGraphQL(rctx, query, v, "")
	if res.HasErrors() {
		t.Error("errors should not exists")
	}
	d := res.Data.(map[string]interface{})
	resData := d["hello"].(string)
	if resData != "world" {
		t.Error("hello response is expected as world")
	}
}
