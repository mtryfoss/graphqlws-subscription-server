package graphqlws_subscription_server

import (
	"github.com/functionalfoundry/graphqlws"
)

type Handler struct {
	manager *graphqlws.SubscriptionManager
	secret  string
}

func NewHandler(manager *graphqlws.SubscriptionManager, secretkey string) *Handler {
	return &Handler{manager: manager, secret: secretkey}
}

func (h *Handler) AuthSecretKey() string {
	return h.secret
}
