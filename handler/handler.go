package graphqlws_subscription_server

import (
	"github.com/graphql-go/graphql"
)

type Handler struct {
	schema   *graphql.Schema
	listener *Listener
}

func NewHandler(listener *Listener) *Handler {
	return &Handler{
		schema:   schema,
		listener: listener,
	}
}
