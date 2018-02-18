package graphqlws_subscription_server

import (
	"github.com/functionalfoundry/graphqlws"
	"github.com/graphql-go/graphql"
)

type Handler struct {
	schema  *graphql.Schema
	manager *graphqlws.SubscriptionManager
}

func NewHandler(schema *graphql.Schema) *Handler {
	manager := graphqlws.NewSubscriptionManager(schema)
	return &Handler{
		schema:  schema,
		manager: &manager,
	}
}
