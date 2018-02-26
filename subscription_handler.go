package gss

import (
	"net/http"

	"github.com/functionalfoundry/graphqlws"
)

func NewSubscriptionHandler(subService *SubscribeService, authCallback graphqlws.AuthenticateFunc) http.Handler {
	return graphqlws.NewHandler(graphqlws.HandlerConfig{
		SubscriptionManager: subService,
		Authenticate:        authCallback,
	})
}
