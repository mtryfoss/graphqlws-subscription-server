package graphqlws_subscription_server

import (
	"net/http"
)

type ChannelNotification struct {
}

type UsersNotification struct {
}

func (r *ChannelNotification) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (r *UsersNotification) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (h *Handler) NewNotifyChannelHandler() http.Handler {
	return &ChannelNotification{}
}

func (h *Handler) NewNotifyUsersHandler() http.Handler {
	return &UsersNotification{}
}
