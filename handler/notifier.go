package graphqlws_subscription_server

import (
	"net/http"
)

type RegistrationResponse struct{}

func (r *RegistrationResponse) ServeHTTP(w http.ResponseWriter, req *http.Request) {

}

func (h *Handler) NewNotifyChannelHandler() http.Handler {
	return &RegistrationResponse{}
}

func (h *Handler) NewNotifyUsersHandler() http.Handler {
	return &RegistrationResponse{}
}
