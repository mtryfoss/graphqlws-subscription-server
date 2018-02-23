package gss

import (
	"net/http"

	"github.com/functionalfoundry/graphqlws"
)

func (h *Handler) NewWebsocketHandler(callback graphqlws.AuthenticateFunc) http.Handler {
	return graphqlws.NewHandler(graphqlws.HandlerConfig{
		SubscriptionManager: h.listener,
		Authenticate:        callback,
	})
}
