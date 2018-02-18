package graphqlws_subscription_server

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/functionalfoundry/graphqlws"
)

func (h *Handler) NewWebsocketHandler() http.Handler {
	return graphqlws.NewHandler(graphqlws.HandlerConfig{
		SubscriptionManager: *h.manager,
		Authenticate: func(tokenstring string) (interface{}, error) {
			token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
				return []byte(h.AuthSecretKey()), nil
			})
			return token.Claims, err
		},
	})
}
